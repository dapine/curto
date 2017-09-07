package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dapine/fng"
	_ "github.com/lib/pq"
)

type Link struct {
	short string
	url   string
}

const ShortLength = 7

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlerLink)

	log.Fatal(http.ListenAndServe(":5000", mux))
}

func handlerLink(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		html := `
		<!DOCTYPE html>
		<html>
			<head>
				<meta charset="utf-8">
				<title>curto - url shortener</title>
				<style>
				/* TODO: Add awesome style */
				</style>
			</head>
			<body>
				<div class="content">
					<h1>curto - url shortener</h1>

					<form action="/new" method="POST">
						<label>Paste your URL here</label>
						<input type="text">
						<input type="submit" value="Short it">
					</form>
				</div>
			</body>
		</html>
		`
		fmt.Fprint(w, html)
	} else {
		db, err := openDB()
		if err != nil {
			fmt.Fprintf(w, "%s", err)
		}

		short := r.URL.Path[1:]

		l, err := QueryLink(db, short)
		if err != nil {
			fmt.Fprint(w, "Not found")
		}

		http.Redirect(w, r, l.url, http.StatusMovedPermanently)
	}
}

func openDB() (*sql.DB, error) {
	userdb := os.Getenv("USERDB")
	dbname := os.Getenv("DBNAME")

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s", userdb, dbname))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewLink(url string) (Link, error) {
	r, err := http.Get(url)
	if err != nil {
		return Link{}, err
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return Link{}, err
	}

	sl, err := fng.GenerateString(b, fng.Charset, ShortLength)
	if err != nil {
		return Link{}, err
	}

	return Link{short: sl, url: url}, nil
}

func InsertLink(db *sql.DB, l Link) error {
	_, err := db.Query(`INSERT INTO link(short, url) VALUES ($1, $2)`, l.short, l.url)
	if err != nil {
		return err
	}

	return nil
}

func QueryLink(db *sql.DB, short string) (Link, error) {
	var s, url string

	rows, err := db.Query("SELECT short, url from link where short = $1", short)
	if err != nil {
		return Link{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&s, &url)
		if err != nil {
			return Link{}, err
		}
	}

	l := Link{short: s, url: url}

	return l, nil
}
