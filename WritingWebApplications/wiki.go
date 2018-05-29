package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

// templates is used to store multiple templates that were parsed, producing 1 Template object.
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// validPath is used to validate page titles submitted by users.
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

//Page represents a wiki page
type Page struct {
	Title string
	Body  []byte
}

//Save is a method on Page struct that will save the page into a text file.
func (page *Page) Save() error {
	filename := page.Title + ".txt"
	return ioutil.WriteFile(filename, page.Body, 0600)
}

//LoadPage loads a text file, reads it and creates a new Page literal from its content.
func LoadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{
		Title: title,
		Body:  body,
	}, nil
}

func main() {
	//this is the routing of endpoint to handler. Like url mappings in web.xml.
	http.HandleFunc("/view/", makeHandler(ViewHandler))
	http.HandleFunc("/edit/", makeHandler(EditHandler))
	http.HandleFunc("/save/", makeHandler(SaveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler is of type http.HandlerFunc (satisfies) because HandlerFunc type is defined with the signature func(ResponseWriter, *Request)
func Handler(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(responseWriter, "Hi there, I love %s!", request.URL.Path[1:])
}

// ViewHandler is a http.HandlerFunc serves a page requested from the "/view/*" endpoint.
// If page doesn't exist, then this handler sends a redirect to the edit page so that the page can be created.
func ViewHandler(responseWriter http.ResponseWriter, request *http.Request, title string) {
	page, err := LoadPage(title)
	if err != nil {
		http.Redirect(responseWriter, request, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(responseWriter, "view", page)
}

// EditHandler is a http.HandlerFunc that serves an edit page.
func EditHandler(responseWriter http.ResponseWriter, request *http.Request, title string) {
	p, err := LoadPage(title)
	if err != nil {
		//if there's error loading page, then just return a page with the title = requested title.
		p = &Page{Title: title}
	}
	renderTemplate(responseWriter, "edit", p)
}

// SaveHandler is a http.HandlerFunc that saves the edit performed on /edit/ page.
func SaveHandler(w http.ResponseWriter, r *http.Request, title string) {
	//get the form value which belongs to the field with name (or key) attribute "body"
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//after save, redirect to view page
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// renderTemplate executes a template that is specified by "tmpl" file name, and is already parsed into the "templates" variable.
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	//execute templates by passing in data (in this case page) that is required to fill up its values.
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// getTitle checks to see if the title submitted by the user is valid. In other words, it can be parsed and compiled according the the validPath regex variable.
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	fmt.Printf("Result of validPath.FindStringSubmatch is : %q\n", m)
	return m[2], nil // The title is the second subexpression.
}

// makeHandler returns a closure http.HandlerFunc that pre-checks the title string in the url. Once the title is validated, it calls the fn function, which are now not satisfying http.HandlerFunc anymore, because of the 3rd string parameter.
// This example of using function literal and closure is a method of abstracting code. DRY principle.
// ** Think of closure functions as a casing of another function.
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	// **important: how to define function using function literal, and then turn it into a closure function by using variables defined outside of the function.
	//CASING
	return func(w http.ResponseWriter, r *http.Request) {
		//CASE CONTENTS
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		//function inside function.
		fn(w, r, m[2])
	}
}
