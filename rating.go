package main

import (
	"fmt"

	"github.com/disgoorg/snowflake/v2"
	"github.com/fluxergo/fluxergo/fluxer"
)

// get steamid and return the rating
func getPremierRatingHandler(author *fluxer.User, message *fluxer.Message, args []string) error {
	profile, err := getLeetifyStats(steamID)
	if err != nil {
		return fmt.Errorf("error getting rating: %w", err)
	}
	premierRating := profile.Ranks.Premier
	fmt.Println("Premier Rating:", premierRating)
	ratingmsg := fluxer.NewMessageCreate().WithContent(fmt.Sprintf("Your premier rating is %s", premierRating))
	_, err = client.Rest.CreateMessage(message.ChannelID, ratingmsg)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}

// updates premier rating role when run
func updatePremierRatingRole(playerstats LeetifyProfile) error {
	premierRating, err := playerstats.Ranks.Premier.Int64()
	if err != nil {
		return fmt.Errorf("error getting premier rating: %w", err)
	}
	roleID := getRoleIDForRating(premierRating)
	fluxerID, err := getFluxerIDFromSteamID(playerstats.SteamID)
	userID := snowflake.MustParse(fluxerID)
	if err != nil {
		return fmt.Errorf("error getting fluxerID: %w", err)
	}
	guildID := snowflake.MustParse("1473790485412413471")
	setPremierRatingRole(guildID, userID, roleID)
	return err

	// needs to check if user has a role and if its the correct role
}

type RatingMap struct {
	Min  int64
	Max  int64
	Role string // fluxer role ID as string
}

var ratings = []RatingMap{
	{Min: 1, Max: 4999, Role: "1474150338416083289"},
	{Min: 5000, Max: 9999, Role: "1474150621913350288"},
	{Min: 10000, Max: 14999, Role: "1474150730499625213"},
	{Min: 15000, Max: 19999, Role: "1474150710530552040"},
	{Min: 20000, Max: 24999, Role: "1474150970833268910"},
	{Min: 25000, Max: 29999, Role: "1474151063175057722"},
	{Min: 30000, Max: 40000, Role: "1474151158499020960"},
}

func getRoleIDForRating(yourRating int64) snowflake.ID {
	for _, rating := range ratings {
		if yourRating >= rating.Min && yourRating <= rating.Max {
			return snowflake.MustParse(rating.Role)
		}
	}
	return snowflake.MustParse("1474240923272941718") // Default role if no match found
}

func setPremierRatingRole(guildID snowflake.ID, userID snowflake.ID, ratingRole snowflake.ID) error {
	err = client.Rest.AddMemberRole(guildID, userID, ratingRole)
	return fmt.Errorf("error adding role: %w", err)
}

// these last two functions can be rewritten
