package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

const meJSONURI string = "https://raw.githubusercontent.com/crmepham/public/master/me.json"
const jsonFETCHERROR string = "Could not fetch JSON!"
const staticCONTENTPATH string = "/static"

// Person encapsulates my public data for use in the template.
type Person struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	JobTitle       string `json:"jobTitle"`
	DateOfBirth    string `json:"dateOfBirth"`
	PersonalEmail  string `json:"personalEmail"`
	PersonalMobile string `json:"personalMobile"`
	ShortBiography string `json:"shortBiography"`
	Links          []Link `json:"links"`
}

// Link encapsulates links to my various online profiles
type Link struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URI         string `json:"uri"`
}

var person = &Person{}

// Redirect to http://chrismepham.co.uk as preferred URL for SEO.
func redirect(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Host, "www") {
		http.Redirect(w, r, "http://chrismepham.co.uk", 301)
	}
}

// The main method.
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc(staticCONTENTPATH+"/{a-z}/{a-z}", staticContentHandler)
	loadPerson()
	log.Fatal(http.ListenAndServe(":8000", r))
}

// The request handlers.
func staticContentHandler(response http.ResponseWriter, request *http.Request) {
	if isCached(response, request) {
		response.WriteHeader(http.StatusNotModified)
		return
	}

	setContentType(response, request)
	path := request.URL.Path
	data, err := ioutil.ReadFile(path[1:len(path)])
	check(err)
	cacheContent(response, request)

	response.Write([]byte(data))
}

func rootHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Link", "http://chrismepham.co.uk")
	redirect(response, request)
	setContentType(response, request)
	temp, err := template.ParseFiles("templates/profile.html")
	check(err)
	temp.Execute(response, person)
}

func jsonHandler(response http.ResponseWriter, request *http.Request) {
	setContentType(response, request)
	unmarshalled, err := json.Marshal(&person)
	check(err)
	response.Write([]byte(string(unmarshalled)))
}

// The utility methods.
func isCached(response http.ResponseWriter, request *http.Request) bool {
	key := request.Header.Get("If-None-Match")
	return key != "" && strings.Contains(key, `"`+request.URL.Path+`"`)
}

func cacheContent(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Etag", `"`+request.URL.Path+`"`)
	response.Header().Set("Cache-Control", "max-age=2592000")
}

func loadPerson() {
	resp, err := http.Get(meJSONURI)
	if err != nil || resp.StatusCode != 200 {
		panic(jsonFETCHERROR)
	}

	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(person)
}

func setContentType(response http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	var suffix = string(path[strings.LastIndex(path, ".")+1 : len(path)])
	var contentType string
	switch suffix {
	case "css":
		contentType = "text/css"
	case "js":
		contentType = "text/javascript"
	case "json":
		contentType = "application/json"
	case "jpg":
		contentType = "image/jpg"
	default:
		contentType = "text/html"
	}
	response.Header().Set("Content-Type", contentType+"; charset=utf-8")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
