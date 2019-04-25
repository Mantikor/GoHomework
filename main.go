// Golang final homework task (c) 2019 Siarhei Straltsou
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"time"
)

const PORT = ":5555"

// TODO use different types for datafields
type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type Account struct {
	ID       int    `json:"id"`
	UserID   int    `json:"accountid"`
	Balance  int    `json:"balance"`
	Opendate string `json:"opendate"`
}

type Transaction struct {
	ID        int    `json:"id"`
	AccountID int    `json:"accountid"`
	DateTime  string `json:"date"`
	Amount    int    `json:"amount"`
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

func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			panic(err)
		}
		users := []User{}
		for rows.Next() {
			p := User{}
			err := rows.Scan(&p.ID, &p.Username, &p.FirstName, &p.LastName, &p.Email, &p.Phone)
			if err != nil {
				fmt.Println(err)
				continue
			}
			users = append(users, p)
		}
		for _, p := range users {
			fmt.Fprintln(w, p.ID, p.Username, p.FirstName, p.LastName, p.Email, p.Phone)

			//w.Header().Set("Content-Type", "application/json")
			//json.NewEncoder(w).Encode(users)
		}
	}
}

func createUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var user User
		_ = json.NewDecoder(r.Body).Decode(&user)
		// TODO change algo for using database
		st, _ := db.Prepare("INSERT INTO users (username, firstname, lastname, email, phone) VALUES (?, ?, ?, ?, ?)")
		st.Exec(user.Username, user.FirstName, user.LastName, user.Email, user.Phone)
		// json.NewEncoder(w).Encode(user)
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Will implement in next version, params: ")
	params := mux.Vars(r)
	json.NewEncoder(w).Encode(params)
	//w.Header().Set("Content-Type", "application/json")
	//params := mux.Vars(r)
	//for index, item := range users {
	//	if item.ID == params["id"] {
	//		users = append(users[:index], users[index+1:]...)
	//		var user User
	//		_ = json.NewDecoder(r.Body).Decode(&user)
	//		user.ID = params["id"]
	//		users = append(users, user)
	//		json.NewEncoder(w).Encode(user)
	//		return
	//	}
	//}
	//// TODO refactor createUser function to assign ID, and remove generaton ID
	//json.NewEncoder(w).Encode(users)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Will implement in next version, params: ")
	params := mux.Vars(r)
	json.NewEncoder(w).Encode(params)
	//w.Header().Set("Content-Type", "application/json")
	//params := mux.Vars(r)
	//for index, item := range users {
	//	if item.ID == params["id"] {
	//		users = append(users[:index], users[index+1:]...)
	//		break
	//	}
	//}
	//json.NewEncoder(w).Encode(users)
}

func getUserAccounts(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Will implement in next version, params: ")
	params := mux.Vars(r)
	json.NewEncoder(w).Encode(params)
}

func createUserAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Will implement in next version, params: ")
	params := mux.Vars(r)
	json.NewEncoder(w).Encode(params)
}

func deleteUserAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Will implement in next version, params: ")
	params := mux.Vars(r)
	json.NewEncoder(w).Encode(params)
}

func getUserAccountTransactions(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Will implement in next version, params: ")
	params := mux.Vars(r)
	json.NewEncoder(w).Encode(params)
}

func main() {
	r := mux.NewRouter()

	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, username TEXT NOT NULL, firstname	TEXT, lastname	TEXT, email	TEXT, phone	TEXT)")
	statement.Exec()
	statement, _ = db.Prepare("CREATE TABLE IF NOT EXISTS accounts (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, user_id INTEGER, balance	INTEGER, opendate	TEXT)")
	statement.Exec()
	statement, _ = db.Prepare("CREATE TABLE IF NOT EXISTS transactions (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, account_id INTEGER, date_time	TEXT, amount INTEGER)")
	statement.Exec()

	users = append(users, User{ID: 1, Username: "user1", FirstName: "John", LastName: "Doe",
		Email: "john.doe@aol.com", Phone: "+190056004"})
	users = append(users, User{ID: 2, Username: "user2", FirstName: "Vasya", LastName: "Pupkin",
		Email: "v.poop@mail.ru", Phone: "+7902586867676"})

	r.HandleFunc("/users", getUsers(db)).Methods("GET")
	r.HandleFunc("/users", createUser(db)).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	// this call use user_id and print all user accounts
	r.HandleFunc("/users/{user_id}/accounts", getUserAccounts).Methods("GET")
	// create account for user_id
	r.HandleFunc("/users/{user_id}/accounts", createUserAccount).Methods("POST")
	// delete acc_id from user_id
	// TODO check account balance == ZERO
	r.HandleFunc("/users/{user_id}/accounts/{acc_id}", deleteUserAccount).Methods("DELETE")
	r.HandleFunc("/users/{user_id}/accounts/{acc_id}", getUserAccountTransactions).Methods("GET")

	if (len(os.Args) > 1) && (os.Args[1] == "--verbose") {
		l := log.New(os.Stdout, "[GoHomeWork-Server] - ", 0)
		l.Printf("Listening 0.0.0.0%s", PORT)
		l.Fatal(http.ListenAndServe(PORT, RequestLogger{h: r, l: l}))
	} else {
		log.Fatal(http.ListenAndServe(PORT, r))
	}
}
