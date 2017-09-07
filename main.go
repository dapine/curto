package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Link struct {
	short string
	url   string
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlerLink)

	log.Fatal(http.ListenAndServe(":5000", mux))
}

func handlerLink(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		fmt.Fprint(w, "Index")
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
