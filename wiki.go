package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

//[]byte rather than string because that is the type
// expected by IO libraries used here
type Page struct {
	Title string
	Body  []byte
}

// This method save as its receiver p, a pointer to the Page
// WriteFile is a standard library function that writes a byte to a file
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// It constructs filename from the title parameter, read file
// contents into a new variable body, return a pointer to a Page
// constructed with proper title and body values
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	t, _ := template.ParseFiles("edit.html")
	t.Execute(w, p)
}

// http.HandleFunc: handle all requests to the web root("/")
//with handler
func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
