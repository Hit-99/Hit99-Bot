package main

import (
	"fmt"
	"strings"

	"github.com/fluxergo/fluxergo/fluxer"
)

// gets fluxerID and steamID and pairs them in db

func linkSteamHandler(author *fluxer.User, message *fluxer.Message, args []string) error {
	msgFluxerID := author.ID
	msgSteamID := args[1]
	if strings.Contains(msgSteamID, "765611") {
		fmt.Println("(debug) steam64 ID: ", msgSteamID)
		// checks if steam64 id
	} else if strings.Contains(msgSteamID, "https://steamcommunity.com/id/") {
		vanityUser := strings.TrimPrefix(msgSteamID, "https://steamcommunity.com/id/")
		vanityUser = strings.TrimSuffix(vanityUser, "/")
		fmt.Println("(debug) steam url: ", vanityUser)
		msgSteamID, err = getSteamID(vanityUser)
		if err != nil {
			return fmt.Errorf("error getting steam id: %w", err)
		}
		// get just the vanity username at the end of link
		// ex. https://steamcommunity.com/id/Hyp3r7/ gets "Hyp3r7"
	} else {
		fmt.Println("(debug) steam misc: ", msgSteamID)
		msgSteamID, err = getSteamID(msgSteamID)
		if err != nil {
			return fmt.Errorf("error getting steam id: %w", err)
		}
		// no link, just username. if not a success, return an error
	}
	fmt.Printf("fluxerID: %s steamID: %s\n", msgFluxerID, msgSteamID)

	if linkContainsID(msgFluxerID) == true {
		fmt.Println("(debug) User already exists")
		alreadyLinkedMsg := fluxer.NewMessageCreate().WithContent(fmt.Sprintf("Your account has already been linked <@!%s>", msgFluxerID))
		_, err = client.Rest.CreateMessage(message.ChannelID, alreadyLinkedMsg)
		if err != nil {
			return fmt.Errorf("error sending message: %w", err)
		}
	} else {
		linkCreateEntry(msgFluxerID, msgSteamID)
		fmt.Println("(debug) User linked")
		linkedMsg := fluxer.NewMessageCreate().WithContent(fmt.Sprintf("Your account has been linked <@!%s>", msgFluxerID))
		_, err = client.Rest.CreateMessage(message.ChannelID, linkedMsg)
		if err != nil {
			return fmt.Errorf("error sending message: %w", err)
		}
		userStats, err := getLeetifyStats(msgSteamID)
		premierRating, _ := userStats.Ranks.Premier.Int64()
		roleID := getRoleIDForRating(premierRating)
		setPremierRatingRole(*message.GuildID, msgFluxerID, roleID)
		return err
	}
	return err

	// needs to check if fluxer id and/or steam id does not already exist
}
