package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

func getId() int {
	id := -1
	for _, p := range people {
		if id < p.Id {
			id = p.Id
		}
	}
	return id + 1
}

func get(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if r.URL.Query().Has("id") {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{
				"errors":[
					{
						"message":"Invalid or no id supplied}
					}
			]}`))
			return
		}

		index := -1
		for i, p := range people {
			if p.Id == id {
				index = i
			}
		}

		if index == -1 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{
				"errors":[
					{"message":"Person with given id not found"}
			]}`))
			return
		}

		p := people[index]
		json.NewEncoder(w).Encode(p)
		w.WriteHeader(http.StatusOK)
		return

	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(people)
}

func post(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"errors":[
				{"message":"Invalid or bad request}
		]}`))
	}

	p := Person{}
	json.Unmarshal(requestBody, &p)

	if p.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"errors":[
				{"message":"Person name is invalid or not supplied"}
		]}`))
		return
	}

	if p.Id == 0 {
		p.Id = getId()
	} else if p.Id != 0 {
		for _, p1 := range people {
			if p1.Id == p.Id {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{
					"errors":[
						{"message":"Person id already present"}
				]}`))
				return
			}
		}
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

	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"errors":[
				{
					"message":"Invalid or no id supplied}
				}
		]}`))
		return
	}

	index := -1
	for i, p := range people {
		if p.Id == id {
			index = i
		}
	}

	if index == -1 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{
			"errors":[
				{"message":"Person with given id not found"}
		]}`))
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"errors":[
				{"message":"Invalid request recieved"}
		]}`))
		return
	}

	var p Person
	json.Unmarshal(requestBody, &p)
	if p.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"errors":[
				{"message":"Person name is invalid or not supplied"}
		]}`))
		return
	}
	people[index].Name = p.Name

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(people[index])
}

func delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"errors":[
				{
					"message":"Invalid or no id supplied}
				}
		]}`))
		return
	}

	index := -1
	for i, p := range people {
		if p.Id == id {
			index = i
		}
	}

	if index == -1 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{
			"errors":[
				{"message":"Person with given id not found"}
		]}`))
		return
	}

	p := people[index]

	people = append(people[:index], people[index+1:]...)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

func main() {
	http.HandleFunc("/", get)
	http.HandleFunc("/create", post)
	http.HandleFunc("/update", put)
	http.HandleFunc("/delete", delete)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
