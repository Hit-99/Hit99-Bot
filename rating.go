package main

import (
	"fmt"

	"github.com/disgoorg/snowflake/v2"
	"github.com/fluxergo/fluxergo/fluxer"
)

// func(caller *fluxer.User, message *fluxer.Message, args []string) error

// !getrating 42384928374982374
// [getrating, 42384928374982374]

func getRatingHandler(author *fluxer.User, message *fluxer.Message, args []string) error {
	rating, err := getRating(args[1])
	if err != nil {
		return fmt.Errorf("error getting rating: %w", err)
	}

	RoleID := getRoleIDForRating(rating)

	setRatingRole(rating, *message.GuildID, author.ID, RoleID)

	ratingmsg := fluxer.NewMessageCreate().WithContent(fmt.Sprintf("Your premier rating is %d. Adding <@&%s> ", rating, RoleID))

	_, err = client.Rest.CreateMessage(message.ChannelID, ratingmsg)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}

func linkRatingHandler(author *fluxer.User, message *fluxer.Message, args []string) error {
	return fmt.Errorf("unimplemented")
}

type RatingMap struct {
	Min  int
	Max  int
	Role string // fluxer role ID as string
}

var ratings = []RatingMap{
	{Min: 0, Max: 4999, Role: "1474150338416083289"},
	{Min: 5000, Max: 9999, Role: "1474150621913350288"},
	{Min: 10000, Max: 14999, Role: "1474150730499625213"},
	{Min: 15000, Max: 19999, Role: "1474150710530552040"},
	{Min: 20000, Max: 24999, Role: "1474150970833268910"},
	{Min: 25000, Max: 29999, Role: "1474151063175057722"},
	{Min: 30000, Max: 40000, Role: "1474151158499020960"},
}

func getRoleIDForRating(yourRating int) snowflake.ID {
	for _, rating := range ratings {
		if yourRating >= rating.Min && yourRating <= rating.Max {
			return snowflake.MustParse(rating.Role)
		}
	}
	return snowflake.MustParse("1474240923272941718") // Default role if no match found
}

func setRatingRole(rating int, guildID snowflake.ID, userID snowflake.ID, ratingRole snowflake.ID) error {

	err = client.Rest.AddMemberRole(guildID, userID, ratingRole)
	return fmt.Errorf("error adding role: %w", err)
}
