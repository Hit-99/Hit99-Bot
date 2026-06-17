package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/disgoorg/snowflake/v2"
	_ "modernc.org/sqlite"
)

var db *sql.DB

type FluxerSteamPair struct {
	FluxerID string
	SteamID  string
}

type DiscordSteamPair struct {
	DiscordID string
	SteamID   string
}

func init() {
	db, err = sql.Open("sqlite", "./hit99.db")
	if err != nil {
		log.Fatal(err)
	}

	// creates fluxerLink table if it doesnt exist
	steamToFluxerLinkSchema := `
	CREATE TABLE IF NOT EXISTS fluxerLink (
		fluxerID TEXT PRIMARY KEY,
		steamID TEXT NOT NULL
	);
	`

	// creates discordLink table if it doesnt exist
	steamToDiscordLinkSchema := `
	CREATE TABLE IF NOT EXISTS discordLink (
		discordID TEXT PRIMARY KEY,
		steamID TEXT NOT NULL
	);
	`

	// creates matches table if it doesnt exist
	steamAndMatchesSchema := `
	CREATE TABLE IF NOT EXISTS matches (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		steamID TEXT NOT NULL,
		matchID TEXT NOT NULL
	);
	`

	_, err = db.Exec(steamToFluxerLinkSchema)
	if err != nil {
		log.Printf("%q: %s\n", err, steamToFluxerLinkSchema)
		log.Println("Unable to create fluxerLink table")
		panic(err)
	}

	_, err = db.Exec(steamToDiscordLinkSchema)
	if err != nil {
		log.Printf("%q: %s\n", err, steamToDiscordLinkSchema)
		log.Println("Unable to create discordLink table")
		panic(err)
	}

	_, err = db.Exec(steamAndMatchesSchema)
	if err != nil {
		log.Printf("%q: %s\n", err, steamAndMatchesSchema)
		log.Println("Unable to create matches table")
		panic(err)
	}
}

// checks if fluxer ID exists in fluxerLink table
func linkContainsFluxerID(fluxerID snowflake.ID) bool {
	var existing string
	err := db.QueryRow("SELECT fluxerID FROM fluxerLink WHERE fluxerID = ?", fluxerID.String()).Scan(&existing)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatal(err)
	}
	return true
}

// checks if discord ID exists in discordLink table
func linkContainsDiscordID(discordID snowflake.ID) bool {
	var existing string
	err := db.QueryRow("SELECT discordID FROM discordLink WHERE discordID = ?", discordID.String()).Scan(&existing)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatal(err)
	}
	return true
}

// creates an entry in the fluxerLink table
func linkCreateEntry(fluxerID snowflake.ID, steamID string) {
	_, err := db.Exec("INSERT INTO fluxerLink (fluxerID, steamID) VALUES (?, ?)", fluxerID.String(), steamID)
	if err != nil {
		log.Fatal(err)
	}
}

// creates an entry in the discordLink table
func dLinkCreateEntry(discordID snowflake.ID, steamID string) {
	_, err := db.Exec("INSERT INTO discordLink (discordID, steamID) VALUES (?, ?)", discordID.String(), steamID)
	if err != nil {
		log.Fatal(err)
	}
}

// get fluxerLink table entries
func linkGetEntries() []FluxerSteamPair {
	rows, err := db.Query("SELECT fluxerID, steamID FROM fluxerLink")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Rows in fluxerLink table:")

	var entries []FluxerSteamPair // create the array outside of the loop

	for rows.Next() {
		var fluxerID, steamID string

		if err := rows.Scan(&fluxerID, &steamID); err != nil { // get your values
			log.Fatal(err)
		}

		entries = append(entries, FluxerSteamPair{FluxerID: fluxerID, SteamID: steamID}) // append to the array

		fmt.Printf("fluxerID=%s steamID=%s\n", fluxerID, steamID)

	}
	return entries
}

// get discordLink table entries
func dLinkGetEntries() []DiscordSteamPair {
	rows, err := db.Query("SELECT discordID, steamID FROM discordLink")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Rows in discordLink table:")

	var entries []DiscordSteamPair // create the array outside of the loop

	for rows.Next() {
		var discordID, steamID string

		if err := rows.Scan(&discordID, &steamID); err != nil { // get your values
			log.Fatal(err)
		}

		entries = append(entries, DiscordSteamPair{DiscordID: discordID, SteamID: steamID}) // append to the array

		fmt.Printf("discordID=%s steamID=%s\n", discordID, steamID)

	}
	return entries
}

// gets all steamIDs in fluxerLink table
func linkGetAllSteamIds() []string {
	var entries []string

	rows, err := db.Query("SELECT steamID FROM fluxerLink")
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

// gets all steamIDs in discordLink table
func dLinkGetAllSteamIds() []string {
	var entries []string

	rows, err := db.Query("SELECT steamID FROM discordLink")
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

func getFluxerIDFromSteamID(steamID string) (string, error) {
	var fluxerID string
	err := db.QueryRow("SELECT fluxerID from fluxerLink where steamID = ?", steamID).Scan(&fluxerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		log.Fatal(err)
	}
	return fluxerID, nil
}

func getDiscordIDFromSteamID(steamID string) (string, error) {
	var discordID string
	err := db.QueryRow("SELECT discordID from discordLink where steamID = ?", steamID).Scan(&discordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		log.Fatal(err)
	}
	return discordID, nil
}

func getSteamIDFromFluxerID(fluxerID string) (string, error) {
	var steamID string
	err := db.QueryRow("SELECT steamID from fluxerLink where fluxerID = ?", fluxerID).Scan(&steamID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		log.Fatal(err)
	}
	return steamID, nil
}

func getSteamIDFromDiscordID(discordID string) (string, error) {
	var steamID string
	err := db.QueryRow("SELECT steamID from discordLink where discordID = ?", discordID).Scan(&steamID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		log.Fatal(err)
	}
	return steamID, nil
}

func ifMatchExistsForUser(steamID string, matchID string) (bool, error) {
	var index string
	err := db.QueryRow("SELECT 1 from matches where steamID = ? AND matchID = ?", steamID, matchID).Scan(&index)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Fatal(err)
	}
	return true, nil
}
