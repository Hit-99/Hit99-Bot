package main

import (
	"fmt"

	"github.com/fluxergo/fluxergo/fluxer"
)

func linkRatingHandler(author *fluxer.User, message *fluxer.Message, args []string) error {
	fluxerID := author.ID
	steamID := args[1]
	fmt.Printf("fluxerID: %s steamID: %s\n", fluxerID, steamID)

	if dbContainsID(fluxerID) != true {
		fmt.Println("User already exists")
	} else {
		dbCreateEntry(fluxerID, steamID)
		fmt.Println("User linked")
	}

	return err
}

/*
user runs !link STEAMID64
the handler will check if the fluxer ID is already in the database
if its not in the db, the user's fluxer ID and steam ID will be paired in the database
*/
