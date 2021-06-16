package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"time"

	nanoid "github.com/matoous/go-nanoid/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/xid"
)

var CREATE_TABLES = `
	CREATE TABLE IF NOT EXISTS oauth_apps(
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		client_id TEXT NOT NULL UNIQUE,
		client_secret TEXT NOT NULL,
		redirect TEXT NOT NULL,
		date DATE NOT NULL
	)
`

func main() {
	db, _ := sql.Open("sqlite3", "database.sqlite")
	_, err := db.Exec(CREATE_TABLES)

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/oauth2/register", func(w http.ResponseWriter, req *http.Request) {
		t, _ := template.ParseFiles("templates/register.html")

		switch req.Method {
		case http.MethodPost:
			name := req.FormValue("name")
			redirect := req.FormValue("redirect")
			client_id := xid.New().String()
			client_secret, _ := nanoid.New()

			stmt, err := db.Prepare(`
				INSERT INTO oauth_apps 
				(name, redirect, client_id, client_secret, date)
				VALUES(?, ?, ?, ?, ?)
			`)

			if err != nil {
				panic(err)
			}

			_, err = stmt.Exec(name, redirect, client_id, client_secret, time.Now().UTC())

			if err != nil {
				panic(err)
			}

			t.Execute(w, map[string]string{
				"Success": "",
			})
		default:
			t.Execute(w, nil)
		}
	})

	http.HandleFunc("/oauth2/apps", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("templates/apps.html")
		rows, _ := db.Query("SELECT name, client_id, date FROM oauth_apps")
		defer rows.Close()

		var app []map[string]string

		for rows.Next() {
			var name string
			var client_id string
			var date time.Time

			rows.Scan(&name, &client_id, &date)

			data := map[string]string{
				"name":      name,
				"client_id": client_id,
				"date":      date.Format("Jan 2, 2006"),
			}

			app = append(app, data)
		}

		t.Execute(w, app)
	})

	err = http.ListenAndServeTLS(":8000", "cert.pem", "key.pem", nil)

	if err != nil {
		fmt.Println(err)
	}
}
