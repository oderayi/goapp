/**
* A simple Go web server
 */

package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I love %s!", r.URL.Path[1:])
}

func run() {
	http.HandleFunc("/", handler)
	/*
	* ListenAndServe always returns error and
	*  blocks until an error occurs.
	* log.Fatal will log the error and end execution.
	* ":9090" means to run the web server on port 9090
	* on any network interface on the server.
	 */
	log.Fatal(http.ListenAndServe(":9090", nil))
}
