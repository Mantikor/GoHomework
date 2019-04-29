// Golang final homework task (c) 2019 Siarhei Straltsou
package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"strconv"
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
// var users []User

func sendTelegramNotify(ChatId, Url, Token, MessageText string) {
	data := []byte(`{"chat_id":"` + ChatId + `", "text":"` + MessageText + `"}`)
	r := bytes.NewReader(data)
	_, _ = http.Post(Url + Token + "/sendMessage", "application/json", r)
}

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
			// FIXME Scan error on column index 3, name "lastname": unsupported Scan, storing driver.Value type <nil> into type *string
			err := rows.Scan(&p.ID, &p.Username, &p.FirstName, &p.LastName, &p.Email, &p.Phone)
			if err != nil {
				fmt.Println(err)
				continue
			}
			users = append(users, p)
		}
		fmt.Fprintln(w, "ID", "Username", "First Name", "Last Name", "E-mail", "Phone number")
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
		// TODO check e-mail and phone format
		st, _ := db.Prepare("INSERT INTO users (username, firstname, lastname, email, phone) VALUES (?, ?, ?, ?, ?)")
		st.Exec(user.Username, user.FirstName, user.LastName, user.Email, user.Phone)
		fmt.Fprintln(w, user.ID, user.Username, user.FirstName, user.LastName, user.Email, user.Phone)
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
	// TODO check account balance == ZERO
	fmt.Fprint(w, "Will implement in next version, params: ")
	params := mux.Vars(r)
	json.NewEncoder(w).Encode(params)
}

func getUserAccountTransactions(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Will implement in next version, params: ")
	params := mux.Vars(r)
	json.NewEncoder(w).Encode(params)

	messageText := "Some message text: " + strconv.Itoa(164646465))
	sendTelegramNotify("11111111", "https://api.telegram.org/bot",
		"777777777:AAAfafuaib95cETvsYFqsqcrwi1K8Z3a4X8", messageText)
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

	r.HandleFunc("/users", getUsers(db)).Methods("GET")
	r.HandleFunc("/users", createUser(db)).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	r.HandleFunc("/users/{user_id}/accounts", getUserAccounts).Methods("GET")
	r.HandleFunc("/users/{user_id}/accounts", createUserAccount).Methods("POST")
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
