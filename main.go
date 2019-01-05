package main

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

const meJSONURI string = "https://raw.githubusercontent.com/crmepham/public/master/me.json"
const jsonFETCHERROR string = "Could not fetch JSON!"

// Person encapsulates my public data for use in the template.
type Person struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DateOfBirth string `json:"dateOfBirth"`
}

var person = &Person{}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	loadPerson()
	log.Fatal(http.ListenAndServe(":8000", r))
}

func loadPerson() {
	resp, err := http.Get(meJSONURI)
	if err != nil || resp.StatusCode != 200 {
		panic(jsonFETCHERROR)
	}

	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(person)
}

func rootHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/html")

	temp, err := template.ParseFiles("templates/profile.html")
	if err != nil {
		response.Write([]byte("Could not load template!"))
		return
	}

	temp.Execute(response, person)
}

func jsonHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	unmarshalled, err := json.Marshal(&person)
	if err != nil {
		panic(jsonFETCHERROR)
	}

	response.Write([]byte(string(unmarshalled)))
}
