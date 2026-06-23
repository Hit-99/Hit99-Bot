package main

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/fluxergo/fluxergo/fluxer"
)

type RatingMap struct {
	Min  int64
	Max  int64
	Role string
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

var dRatings = []RatingMap{
	{Min: 1, Max: 4999, Role: "1474815503787888660"},
	{Min: 5000, Max: 9999, Role: "1474815590995853444"},
	{Min: 10000, Max: 14999, Role: "1474815692384501834"},
	{Min: 15000, Max: 19999, Role: "1474815818297643089"},
	{Min: 20000, Max: 24999, Role: "1474815904788250776"},
	{Min: 25000, Max: 29999, Role: "1474816063425478706"},
	{Min: 30000, Max: 40000, Role: "1474816175283241133"},
}

// get steamid and return the rating (fluxer)
func getPremierRatingHandler(author *fluxer.User, message *fluxer.Message, args []string) error {
	steamID, err := getSteamIDFromFluxerID(author.ID.String())
	if err != nil {
		return fmt.Errorf("error getting steamID: %w", err)
	}

	profile, err := getLeetifyStats(steamID)
	if err != nil {
		return fmt.Errorf("error getting rating: %w", err)
	}
	premierRating := profile.Ranks.Premier
	ratingmsg := fluxer.NewMessageCreate().WithContent(fmt.Sprintf("Your premier rating is %s", premierRating))
	_, err = fClient.Rest.CreateMessage(message.ChannelID, ratingmsg)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}

// get steamid and return the rating (discord)
func dGetPremierRatingHandler(author *discord.User, message *discord.Message, args []string) error {
	steamID, err := getSteamIDFromDiscordID(author.ID.String())
	if err != nil {
		return fmt.Errorf("error getting steamID: %w", err)
	}

	profile, err := getLeetifyStats(steamID)
	if err != nil {
		return fmt.Errorf("error getting rating: %w", err)
	}
	premierRating := profile.Ranks.Premier
	ratingmsg := discord.NewMessageCreate().WithContent(fmt.Sprintf("Your premier rating is %s", premierRating))
	_, err = dClient.Rest.CreateMessage(message.ChannelID, ratingmsg)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}

// update users premier rating role based on current leetify stats (fluxer)
func updatePremierRatingRole(playerstats LeetifyProfile) error {
	var premierRating int64

	// returns nil if no leetify account
	if playerstats.SteamID == "" {
		return nil
	}

	// check for unrated players
	if playerstats.Ranks.Premier == "" {
		premierRating = 0
	} else {
		premierRating, err = playerstats.Ranks.Premier.Int64()
	}
	if err != nil {
		return fmt.Errorf("error getting premier rating: %w", err)
	}

	roleID := getRoleIDForRating(premierRating)
	fluxerID, err := getFluxerIDFromSteamID(playerstats.SteamID)
	if err != nil {
		return fmt.Errorf("error getting fluxerID: %w", err)
	}
	if fluxerID == "" {
		return nil
	}

	userID := snowflake.MustParse(fluxerID)
	guildID := snowflake.MustParse("1473790485412413471")

	member, err := fClient.Rest.GetMember(guildID, userID)
	if err != nil {
		return err
	}

	// check if user already has the role
	hasCorrectRole := false
	for _, ownedRoleID := range member.RoleIDs {
		if ownedRoleID == roleID {
			hasCorrectRole = true
			break
		}
	}

	// check all user roles and remove any that the user should not have
	for _, r := range member.RoleIDs {
		for _, rating := range ratings {
			if r.String() == rating.Role && rating.Role != roleID.String() {
				ratingRole := snowflake.MustParse(rating.Role)
				err := removePremierRatingRole(guildID, userID, ratingRole)
				if err != nil {
					fmt.Println("failed removing role:", err)
				}
				time.Sleep(250 * time.Millisecond)
			}
		}
	}

	// set the correct premier rating role
	if !hasCorrectRole {
		err = setPremierRatingRole(guildID, userID, roleID)
		if err != nil {
			return err
		}
	}

	return nil
}

// update users premier rating role based on current leetify stats (discord)
func dUpdatePremierRatingRole(playerstats LeetifyProfile) error {
	var premierRating int64

	// returns nil if no leetify account
	if playerstats.SteamID == "" {
		return nil
	}

	// check for unrated players
	if playerstats.Ranks.Premier == "" {
		premierRating = 0
	} else {
		premierRating, err = playerstats.Ranks.Premier.Int64()
	}
	if err != nil {
		return fmt.Errorf("error getting premier rating: %w", err)
	}

	roleID := dGetRoleIDForRating(premierRating)
	discordID, err := getDiscordIDFromSteamID(playerstats.SteamID)
	if err != nil {
		return fmt.Errorf("error getting discordID: %w", err)
	}
	if discordID == "" {
		return nil
	}

	userID := snowflake.MustParse(discordID)
	guildID := snowflake.MustParse("708918939308130345")

	member, err := dClient.Rest.GetMember(guildID, userID)
	if err != nil {
		return err
	}

	// check if user already has the role
	hasCorrectRole := false
	for _, ownedRoleID := range member.RoleIDs {
		if ownedRoleID == roleID {
			hasCorrectRole = true
			break
		}
	}

	// check all user roles and remove any that the user should not have
	for _, r := range member.RoleIDs {
		for _, rating := range dRatings {
			if r.String() == rating.Role && rating.Role != roleID.String() {
				ratingRole := snowflake.MustParse(rating.Role)
				err := dRemovePremierRatingRole(guildID, userID, ratingRole)
				if err != nil {
					fmt.Println("failed removing role:", err)
				}
				time.Sleep(250 * time.Millisecond)
			}
		}
	}

	// set the correct premier rating role
	if !hasCorrectRole {
		err = dSetPremierRatingRole(guildID, userID, roleID)
		if err != nil {
			return err
		}
	}

	return nil
}

// get fluxer role ID based on leetify rating
func getRoleIDForRating(yourRating int64) snowflake.ID {
	for _, rating := range ratings {
		if yourRating >= rating.Min && yourRating <= rating.Max {
			return snowflake.MustParse(rating.Role)
		}
	}
	return snowflake.MustParse("1474240923272941718") // Default role if no match found
}

// get discord role ID based on leetify rating
func dGetRoleIDForRating(yourRating int64) snowflake.ID {
	for _, rating := range dRatings {
		if yourRating >= rating.Min && yourRating <= rating.Max {
			return snowflake.MustParse(rating.Role)
		}
	}
	return snowflake.MustParse("1474815393074774178") // Default role if no match found
}

// set specified role for user (fluxer)
func setPremierRatingRole(guildID, userID, ratingRole snowflake.ID) error {
	return fClient.Rest.AddMemberRole(guildID, userID, ratingRole)
}

// set specified role for user (discord)
func dSetPremierRatingRole(guildID, userID, ratingRole snowflake.ID) error {
	return dClient.Rest.AddMemberRole(guildID, userID, ratingRole)
}

// removed specified role from user (fluxer)
func removePremierRatingRole(guildID, userID, role snowflake.ID) error {
	return fClient.Rest.RemoveMemberRole(guildID, userID, role)
}

// removed specified role from user (discord)
func dRemovePremierRatingRole(guildID, userID, role snowflake.ID) error {
	return dClient.Rest.RemoveMemberRole(guildID, userID, role)
}
