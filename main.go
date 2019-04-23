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

const PORT = ":5555"

// TODO use different types for datafields
type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type Account struct {
	ID        string `json:"id"`
	AccountID string `json:"accountid"`
	Balance   string `json:"balance"`
	Opendate  string `json:"opendate"`
}

type Transaction struct {
	AccountID string `json:"id"`
	Date      string `json:"date"`
	Amount    string `json:"amount"`
}

type RequestLogger struct {
	h http.Handler
	l *log.Logger
}

// TODO for test only, will remove after DB connect
var users []User

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

func getUserAccountTransactions(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "will implement in next version!")
}

func main() {
	r := mux.NewRouter()

	users = append(users, User{ID: "1", Username: "user1", FirstName: "John", LastName: "Doe",
		Email: "john.doe@aol.com", Phone: "+190056004"})
	users = append(users, User{ID: "2", Username: "user2", FirstName: "Vasya", LastName: "Pupkin",
		Email: "v.poop@mail.ru", Phone: "+7902586867676"})

	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	// this call use user_id and print all user accounts
	r.HandleFunc("/accounts/{user_id}", getUserAccounts).Methods("GET")
	// create account for user_id
	r.HandleFunc("/accounts/{user_id}", createUserAccount).Methods("POST")
	// delete acc_id from user_id
	// TODO check account balance == ZERO
	r.HandleFunc("/accounts/{user_id}/{acc_id}", deleteUserAccount).Methods("DELETE")
	r.HandleFunc("/accounts/{user_id}/{acc_id}", getUserAccountTransactions).Methods("GET")

	fmt.Println(os.Args)
	if (len(os.Args) > 1) && (os.Args[1] == "--verbose") {
		l := log.New(os.Stdout, "[GoHomeWork-Server] - ", 0)
		l.Printf("Listening 0.0.0.0%s", PORT)
		l.Fatal(http.ListenAndServe(PORT, RequestLogger{h: r, l: l}))
	} else {
		log.Fatal(http.ListenAndServe(PORT, r))
	}
}
