package main

import (
	"fmt"

	"github.com/fluxergo/fluxergo/fluxer"
)

// gets fluxerID and steamID and pairs them in db

func linkRatingHandler(author *fluxer.User, message *fluxer.Message, args []string) error {
	msgFluxerID := author.ID
	msgSteamID := args[1]
	fmt.Printf("fluxerID: %s steamID: %s\n", msgFluxerID, msgSteamID)

	if linkContainsID(msgFluxerID) == true {
		fmt.Println("User already exists")
	} else {
		linkCreateEntry(msgFluxerID, msgSteamID)
		fmt.Println("User linked")
	}
	return err
	// needs to check if fluxer id and/or steam id does not already exist
}
