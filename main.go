// Golang final homework task (c) 2019 Siarhei Straltsou
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Accounts  []Account
}

type Account struct {
	ID       string `json:"id"`
	Balance  string `json:"balance"`
	Opendate string `json:"opendate"`
}

type RequestLogger struct {
	h http.Handler
	l *log.Logger
}

var users []User

var accounts []Account

func (rl RequestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	rl.l.Printf("Started %s %s", r.Method, r.URL.Path)
	rl.h.ServeHTTP(w, r)
	rl.l.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	// TODO change algo for using database
	user.ID = strconv.Itoa(rand.Intn(1000000))
	users = append(users, user)
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range users {
		if item.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			var user User
			_ = json.NewDecoder(r.Body).Decode(&user)
			user.ID = params["id"]
			users = append(users, user)
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	// TODO refactor createUser function to assign ID, and remove generaton ID
	createUser(w, r)
	json.NewEncoder(w).Encode(users)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range users {
		if item.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(users)
}

func getUserAccounts(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "will implement in next version!")
}

func createUserAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "will implement in next version!")
}

func deleteUserAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "will implement in next version!")
}

func main() {
	r := mux.NewRouter()

	users = append(users, User{ID: "1", Username: "user1", FirstName: "John", LastName: "Doe",
		Email: "john.doe@aol.com", Phone: "+190056004", Accounts: []Account{}})
	users = append(users, User{ID: "2", Username: "user2", FirstName: "Vasya", LastName: "Pupkin",
		Email: "v.poop@mail.ru", Phone: "+7902586867676", Accounts: []Account{}})

	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	r.HandleFunc("/accounts/{id}", getUserAccounts).Methods("GET")
	r.HandleFunc("/accounts/{id}", createUserAccount).Methods("POST")
	r.HandleFunc("/accounts/{id}", deleteUserAccount).Methods("DELETE")

	//log.Fatal(http.ListenAndServe(":5555", r))
	// TODO implement -verbose flag
	l := log.New(os.Stdout, "[GoHomeWork-Server] ", 0)
	l.Printf("Listening 0.0.0.0%s", ":5555")
	l.Fatal(http.ListenAndServe(":5555", RequestLogger{h: r, l: l}))
}
