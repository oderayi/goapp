package main

import (
	"fmt"
	"io/ioutil"
)

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
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

/**
*	Declare load method on Page
*
* ReadFile returns two values: file content and error
* Returns a pointer to Page literal
 */
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

/**
* main function
 */
func main() {
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
}
