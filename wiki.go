package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var templateDir = "./tmpl"
var dataDir = "./data"

/**
* Parse and cache all templates for performance
* template.Must is used to panic if err is not nil.
* If that happens, the program is exited.
*
* template.ParseFiles parses all the templates files into an array
* of templates indexed by the template files base names
 */
var templates = template.Must(template.ParseFiles(templateDir+"/edit.html", templateDir+"/view.html"))

/**
* regular expression for valid path
* MustCompile parses and compiles the regular expression.
* Will panic if compile fails.
 */
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

/**
* Regex for searching for [PageName] for inter-page links
 */
var interPage = regexp.MustCompile(`\[([a-zA-Z0-9]+)\]`)

/*Page structure */
type Page struct {
	Title string
	Body  []byte
}

/* Declare a save() method on Page struct
*
* (p *Page) is a receiver to the function save
* error is the return type. However, if the write is
* successful, save will return `nil`
 */
func (p *Page) save() error {
	filename := dataDir + "/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

/**
*	Declare load method on Page
*
* ReadFile returns two values: file content and error
* Returns a pointer to Page literal
 */
func loadPage(title string) (*Page, error) {
	filename := dataDir + "/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

/**
* Replace inter-page links variables [PageName]
* with html links
 */
func replaceInterPageLinks(body []byte) []byte {
	body = interPage.ReplaceAllFunc(body, func(s []byte) []byte {
		m := string(s[1 : len(s)-1])
		return []byte("<a href=\"/view/" + m + "\"> " + m + "</a>")
	})
	return body
}

/**
* http handler for frontpage
 */
func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
}

/**
* http handler for viewing wiki
 */
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	nP := &Page{Title: p.Title, Body: replaceInterPageLinks(p.Body)}
	renderTemplate(w, "view", nP)
}

/**
* http handler for editing wiki
 */
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

/**
* http handler for saving wiki
 */
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

/**
* A function that creates HandlerFunc functions
* Takes any of the handler functions (view, edit, save)
* Returns a closure which is compatible with http.HandleFunc
 * The variable fn is enclosed by the closure
*/
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract Title
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		// call handler
		fn(w, r, m[2])
	}
}

/**
* Extract title from URL, safely
 */
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression
}

/**
* Render template
 */
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/**
* main function
 */
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":9090", nil))
}
