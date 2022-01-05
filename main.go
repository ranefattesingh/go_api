package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Person struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

var people = []Person{
	{
		Id:   1,
		Name: "person1",
	},
	{
		Id:   2,
		Name: "person2",
	},
}

func get(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}

func post(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("ERROR", err)
	}

	p := Person{}
	json.Unmarshal(requestBody, &p)

	if p.Id == 0 {
		p.Id = len(people) + 1
	}
	people = append(people, p)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func put(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	query := &r.URL
	fmt.Println(query)
	//TODO many things
	w.WriteHeader(http.StatusOK)
}

func delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	query := &r.URL
	fmt.Println(query)
	// TODO many things
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", get)
	http.HandleFunc("/create", post)
	http.HandleFunc("/update/{id}", put)
	http.HandleFunc("/delete/{id}", delete)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
