// Golang final homework task (c) 2019 Siarhei Straltsou
package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	//	Account   accounts []Account
}

//type Account struct {
//	Balance  string `json:"balance"`
//	Opendate string `json:"opendate"`
//}

var users []User

//var accounts []Account

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	// ToDo: change algo for database
	user.ID = strconv.Itoa(rand.Intn(1000000))
	users = append(users, user)
	json.NewEncoder(w).Encode(user)
}

func main() {
	r := mux.NewRouter()

	users = append(users, User{ID: "1", Username: "user1", FirstName: "John", LastName: "Doe",
		Email: "john.doe@aol.com", Phone: "+190056004"})
	users = append(users, User{ID: "2", Username: "user2", FirstName: "Vasya", LastName: "Pupkin",
		Email: "v.poop@mail.ru", Phone: "+7902586867676"})

	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":5555", r))
}