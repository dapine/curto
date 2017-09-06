package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Link struct {
	short string
	url   string
}

func main() {
	userdb := os.Getenv("USERDB")
	dbname := os.Getenv("DBNAME")

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s", userdb, dbname))
	if err != nil {
		log.Fatal(err)
	}

	l, err := QueryLink(db, "abc")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(l)
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
		err = rows.Scan(&s, &url)
		if err != nil {
			return Link{}, err
		}
	}

	l := Link{short: s, url: url}

	return l, nil
}
