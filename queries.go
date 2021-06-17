package main

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

var CREATE_AUTH_APP = `
	INSERT INTO oauth_apps 
	(name, redirect, client_id, client_secret, date)
	VALUES(?, ?, ?, ?, ?)
`

var GET_APPS = "SELECT name, client_id, date FROM oauth_apps"

var GET_APP = `
	SELECT name, client_id, client_secret, redirect, date FROM oauth_apps
	WHERE client_id = ?
`
