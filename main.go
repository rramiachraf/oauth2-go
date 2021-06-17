package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	nanoid "github.com/matoous/go-nanoid/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/xid"
)

var db, _ = sql.Open("sqlite3", "database.sqlite")

func main() {
	defer db.Close()
	db.Exec(CREATE_TABLES)

	http.HandleFunc("/oauth2/register", registerApp)
	http.HandleFunc("/oauth2/apps", myApps)
	http.HandleFunc("/oauth2/app/", viewApp)

	err := http.ListenAndServeTLS(":8000", "cert.pem", "key.pem", nil)

	if err != nil {
		fmt.Println(err)
	}
}

func registerApp(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("templates/register.tmpl", "templates/navbar.tmpl")

	switch req.Method {
	case http.MethodPost:
		name := req.FormValue("name")
		redirect := req.FormValue("redirect")
		client_id := xid.New().String()
		client_secret, _ := nanoid.New()

		stmt, err := db.Prepare(CREATE_AUTH_APP)

		if err != nil {
			panic(err)
		}

		_, err = stmt.Exec(name, redirect, client_id, client_secret, time.Now().UTC())

		if err != nil {
			panic(err)
		}

		http.Redirect(w, req, "/oauth2/apps", http.StatusMovedPermanently)
	default:
		t.Execute(w, nil)
	}
}

func myApps(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/apps.tmpl", "templates/navbar.tmpl")
	rows, _ := db.Query(GET_APPS)
	defer rows.Close()

	var Apps []map[string]string

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

		Apps = append(Apps, data)
	}

	t.Execute(w, Apps)
}

func viewApp(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("templates/app.tmpl", "templates/navbar.tmpl")
	// Too lazy to install gorilla/mux :/
	url := strings.Split(req.RequestURI, "/")
	id := url[len(url)-1]

	row := db.QueryRow(GET_APP, id)

	var name string
	var client_id string
	var client_secret string
	var redirect string
	var date time.Time

	row.Scan(&name, &client_id, &client_secret, &redirect, &date)

	t.Execute(w, map[string]string{
		"name":          name,
		"client_id":     client_id,
		"client_secret": client_secret,
		"redirect":      redirect,
		"date":          date.Format("Monday January 2, 2006"),
	})
}
