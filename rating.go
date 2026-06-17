package main

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/fluxergo/fluxergo/fluxer"
)

// get steamid and return the rating
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
	fmt.Println("Premier Rating:", premierRating)
	ratingmsg := fluxer.NewMessageCreate().WithContent(fmt.Sprintf("Your premier rating is %s", premierRating))
	_, err = fClient.Rest.CreateMessage(message.ChannelID, ratingmsg)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
	// if unrated, message is empty
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
	fmt.Println("Premier Rating:", premierRating)
	ratingmsg := discord.NewMessageCreate().WithContent(fmt.Sprintf("Your premier rating is %s", premierRating))
	_, err = dClient.Rest.CreateMessage(message.ChannelID, ratingmsg)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
	// if unrated, message is empty
}

// updates premier rating role when run
func updatePremierRatingRole(playerstats LeetifyProfile) error {
	premierRating, err := playerstats.Ranks.Premier.Int64()
	if err != nil {
		return fmt.Errorf("error getting premier rating: %w", err)
	}

	roleID := getRoleIDForRating(premierRating)
	fluxerID, err := getFluxerIDFromSteamID(playerstats.SteamID)

	if err != nil {
		return fmt.Errorf("error getting fluxerID: %w", err)
	}

	userID := snowflake.MustParse(fluxerID)
	guildID := snowflake.MustParse("1473790485412413471")

	member, err := fClient.Rest.GetMember(guildID, userID)
	if err != nil {
		return err
	}

	for _, ownedRoleID := range member.RoleIDs {
		if ownedRoleID == roleID {
			return nil
		}
	}

	for _, r := range member.RoleIDs {
		for _, rating := range ratings {
			ratingRole := snowflake.MustParse(rating.Role)
			if r == ratingRole {
				err := removePremierRatingRole(guildID, userID, ratingRole)
				if err != nil {
					fmt.Println("Failed removing role:", err)
				}
				time.Sleep(250 * time.Millisecond)
			}
		}
	}

	err = setPremierRatingRole(guildID, userID, roleID)

	if err != nil {
		return err
	}
	return nil
}

func dUpdatePremierRatingRole(playerstats LeetifyProfile) error {
	premierRating, err := playerstats.Ranks.Premier.Int64()
	if err != nil {
		return fmt.Errorf("error getting premier rating: %w", err)
	}

	roleID := dGetRoleIDForRating(premierRating)
	discordID, err := getDiscordIDFromSteamID(playerstats.SteamID)

	if err != nil {
		return fmt.Errorf("error getting discordID: %w", err)
	}

	userID := snowflake.MustParse(discordID)
	guildID := snowflake.MustParse("708918939308130345")

	member, err := dClient.Rest.GetMember(guildID, userID)
	if err != nil {
		return err
	}

	for _, ownedRoleID := range member.RoleIDs {
		if ownedRoleID == roleID {
			return nil
		}
	}

	for _, r := range member.RoleIDs {
		for _, rating := range ratings {
			ratingRole := snowflake.MustParse(rating.Role)
			if r == ratingRole {
				err := dRemovePremierRatingRole(guildID, userID, ratingRole)
				if err != nil {
					fmt.Println("Failed removing role:", err)
				}
				time.Sleep(250 * time.Millisecond)
			}
		}
	}

	err = dSetPremierRatingRole(guildID, userID, roleID)

	if err != nil {
		return err
	}
	return nil
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

var dRatings = []RatingMap{
	{Min: 1, Max: 4999, Role: "1474815503787888660"},
	{Min: 5000, Max: 9999, Role: "1474815590995853444"},
	{Min: 10000, Max: 14999, Role: "1474815692384501834"},
	{Min: 15000, Max: 19999, Role: "1474815818297643089"},
	{Min: 20000, Max: 24999, Role: "1474815904788250776"},
	{Min: 25000, Max: 29999, Role: "1474816063425478706"},
	{Min: 30000, Max: 40000, Role: "1474816175283241133"},
}

func getRoleIDForRating(yourRating int64) snowflake.ID {
	for _, rating := range ratings {
		if yourRating >= rating.Min && yourRating <= rating.Max {
			return snowflake.MustParse(rating.Role)
		}
	}
	return snowflake.MustParse("1474240923272941718") // Default role if no match found
}

func dGetRoleIDForRating(yourRating int64) snowflake.ID {
	for _, rating := range dRatings {
		if yourRating >= rating.Min && yourRating <= rating.Max {
			return snowflake.MustParse(rating.Role)
		}
	}
	return snowflake.MustParse("1474815393074774178") // Default role if no match found
}

func setPremierRatingRole(guildID, userID, ratingRole snowflake.ID) error {
	return fClient.Rest.AddMemberRole(guildID, userID, ratingRole)
}

func dSetPremierRatingRole(guildID, userID, ratingRole snowflake.ID) error {
	return dClient.Rest.AddMemberRole(guildID, userID, ratingRole)
}

func removePremierRatingRole(guildID, userID, role snowflake.ID) error {
	return fClient.Rest.RemoveMemberRole(guildID, userID, role)
}

func dRemovePremierRatingRole(guildID, userID, role snowflake.ID) error {
	return dClient.Rest.RemoveMemberRole(guildID, userID, role)
}
