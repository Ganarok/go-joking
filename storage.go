package main

import (
	"database/sql"
	"io"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func handleDatabaseFiles() error {
	errCheck, isCreated := createSqliteFile()

	if errCheck != nil {
		return errCheck
	}

	db, err := sql.Open("sqlite3", "./db/storage")

	if err != nil {
		return err
	}

	errConn := databaseConnection(isCreated, db)

	if errConn != nil {
		return errConn
	}

	return nil
}

func createSqliteFile() (error, bool) {
	_, err := os.Stat("./db/storage.sqlite")

	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create("./db/storage.sqlite")

			if err != nil {
				return err, false
			}

			defer f.Close()

			return nil, true
		} else {
			return err, false
		}
	}

	return nil, false
}

func seedDatabase(db *sql.DB) error {
	f, err := os.Open("./db/storage.sqlite")

	if err != nil {
		return err
	}

	data, err := io.ReadAll(f)
	log.Println("Data length: ", len(data))

	if err != nil {
		return err
	}

	if len(data) == 0 {
		log.Println("Replacing empty sqlite with default sqlite.")
		fDefault, err := os.Open("./db/default.sqlite")

		if err != nil {
			return err
		}

		_, err = io.Copy(f, fDefault)

		if err != nil {
			return err
		}

		defer fDefault.Close()
	}

	defer f.Close()

	sqlQuery := string(data)

	_, err = db.Exec(sqlQuery)

	if err != nil {
		return err
	}

	return nil
}

func databaseConnection(isCreated bool, db *sql.DB) error {
	if isCreated {
		seedDatabase(db)
	}

	DB = db

	return nil
}

func getJoke() (Joke, error) {
	row, err := DB.Query("SELECT question, answer FROM jokes ORDER BY RANDOM() LIMIT 1")

	if err != nil {
		return Joke{}, err
	}

	defer row.Close()

	var joke Joke

	for row.Next() {
		err = row.Scan(&joke.question, &joke.answer)

		log.Println(joke)

		if err != nil {
			return Joke{}, err
		}
	}

	err = row.Err()

	if err != nil {
		return Joke{}, err
	}

	return joke, nil
}
