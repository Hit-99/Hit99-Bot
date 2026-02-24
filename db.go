package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/disgoorg/snowflake/v2"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func init() {
	db, err = sql.Open("sqlite", "./hit99.db")
	if err != nil {
		log.Fatal(err)
	}

	// creates tables if they dont exist
	steamToFluxerLinkSchema := `
	CREATE TABLE IF NOT EXISTS link (
		fluxerID TEXT PRIMARY KEY,
		steamID TEXT NOT NULL
	);
	`

	_, err = db.Exec(steamToFluxerLinkSchema)
	if err != nil {
		log.Printf("%q: %s\n", err, steamToFluxerLinkSchema)
		log.Println("Unable to create link table")
		panic(err)
	}
}

// checks if fluxer ID exists in db
func dbContainsID(fluxerID snowflake.ID) bool {
	var existing string
	err := db.QueryRow("SELECT fluxerID FROM link WHERE fluxerID = ?", fluxerID.String()).Scan(&existing)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatal(err)
	}
	fmt.Println(db.Query("SELECT * FROM link"))
	return true
}

// creates an entry in the db
func dbCreateEntry(fluxerID snowflake.ID, steamID string) {
	_, err := db.Exec("INSERT INTO link (fluxerID, steamID) VALUES (?, ?)", fluxerID.String(), steamID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(db.Query("SELECT * FROM link"))
}

// something here does not work, but i dont know what yet
