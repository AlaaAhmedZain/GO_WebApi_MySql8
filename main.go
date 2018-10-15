package main

import (
	"./src/github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// The person Type (more like an object)
type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

type City struct{
	ID int
	Name string
	CountryCode string
	District string
	Population int
}

var people []Person
var citiy []City

// Display all from the people var
func GetPeople(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}

// Display a single data
func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

// create a new item
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

// Delete an item
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}


// Display from the city var
func GetCity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range citiy {
		if strconv.Itoa(item.ID) == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	json.NewEncoder(w).Encode(citiy)
}

// create a new item
func CreateCity(w http.ResponseWriter, r *http.Request) {

	var c City
	_ = json.NewDecoder(r.Body).Decode(&c)

	db, err := sql.Open("mysql", "root:PASS@WORD123@tcp(127.0.0.1:3306)/world")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// perform a db.Query insert
	insert, err := db.Query("INSERT INTO City VALUES ( "+strconv.Itoa(c.ID)+",'"+c.Name+"'"+",'"+c.CountryCode+"'"+ ",'"+c.District+"'"+ ","+strconv.Itoa(c.Population)+  " )")

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()

	citiy = append(citiy, c)
	json.NewEncoder(w).Encode(people)
}

// Delete an item
func DeleteCity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	db, err := sql.Open("mysql", "root:PASS@WORD123@tcp(127.0.0.1:3306)/world")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// perform a db.Query insert
	del, err := db.Query("delete from City where ID=  "+params["id"])

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer del.Close()


}


// main function to boot up everything
func main() {
	router := mux.NewRouter()

	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")

	db, err := sql.Open("mysql", "root:PASS@WORD123@tcp(127.0.0.1:3306)/world")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Execute the query
	results, err := db.Query("SELECT * FROM city")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var c City
		// for each row, scan the result into our tag composite object
		err = results.Scan(&c.ID, &c.Name,&c.CountryCode,&c.District,&c.Population)

		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute

		citiy=append(citiy,c)

	}
	router.HandleFunc("/city", GetCity).Methods("GET")
	router.HandleFunc("/city/{id}", GetCity).Methods("GET")
	router.HandleFunc("/city", CreateCity).Methods("POST")
	router.HandleFunc("/city/{id}", DeleteCity).Methods("DELETE")



	log.Fatal(http.ListenAndServe(":8000", router))
}