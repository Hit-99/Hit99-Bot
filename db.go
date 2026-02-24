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

	// creates link table if it doesnt exist
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

	// creates matches table if it doesnt exist
	steamAndMatchesSchema := `
	CREATE TABLE IF NOT EXISTS matches (
		ID UUID PRIMARY KEY,
		steamID TEXT NOT NULL,
		matchID TEXT NOT NULL
	);
	`

	_, err = db.Exec(steamAndMatchesSchema)
	if err != nil {
		log.Printf("%q: %s\n", err, steamAndMatchesSchema)
		log.Println("Unable to create matches table")
		panic(err)
	}
}

// checks if fluxer ID exists in link table
func linkContainsID(fluxerID snowflake.ID) bool {
	var existing string
	err := db.QueryRow("SELECT fluxerID FROM link WHERE fluxerID = ?", fluxerID.String()).Scan(&existing)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatal(err)
	}
	return true
}

// creates an entry in the link table
func linkCreateEntry(fluxerID snowflake.ID, steamID string) {
	_, err := db.Exec("INSERT INTO link (fluxerID, steamID) VALUES (?, ?)", fluxerID.String(), steamID)
	if err != nil {
		log.Fatal(err)
	}
}

type FluxerSteamPair struct {
	FluxerID string
	SteamID  string
}

// get link table entries
func linkGetEntries() []FluxerSteamPair {
	rows, err := db.Query("SELECT fluxerID, steamID FROM link")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Rows in link table:")

	var entries []FluxerSteamPair // create the array outside of the loop

	for rows.Next() {
		var fluxerID, steamID string

		if err := rows.Scan(&fluxerID, &steamID); err != nil { // get your values
			log.Fatal(err)
		}

		entries = append(entries, FluxerSteamPair{FluxerID: fluxerID, SteamID: steamID}) // append to the array

		fmt.Printf("fluxerID=%s steamID=%s\n", fluxerID, steamID)

	}
	return entries // return the array
}

// gets all steamIDs in link table
func linkGetAllSteamIds() []string {
	var entries []string

	rows, err := db.Query("SELECT steamID FROM link")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var steamID string

		if err := rows.Scan(&steamID); err != nil { // get your values
			log.Fatal(err)
		}

		entries = append(entries, steamID) // append to the array
	}

	return entries
}

// checks if match ID exists in match table
func matchesContainsMatchID(matchID string) bool {
	var existing string
	err := db.QueryRow("SELECT matchID FROM matches WHERE matchID = ?", matchID).Scan(&existing)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatal(err)
	}
	return true
}

// checks if steam ID exists in match table
func matchesContainsSteamID(steamID string) bool {
	var existing string
	err := db.QueryRow("SELECT steamID FROM matches WHERE steamID = ?", steamID).Scan(&existing)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatal(err)
	}
	return true
}

// creates an entry in the match table
func matchesCreateEntry(steamID string, matchID string) {
	_, err := db.Exec("INSERT INTO matches (steamID, matchID) VALUES (?, ?)", steamID, matchID)
	if err != nil {
		log.Fatal(err)
	}
}
