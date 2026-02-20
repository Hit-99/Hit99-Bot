package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/disgoorg/snowflake/v2"
	"github.com/fluxergo/fluxergo/fluxer"
)

type Profile struct {
	Ranks struct {
		Premier int `json:"premier"`
	} `json:"ranks"`
}

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

func setRatingRole(rating int, guildID snowflake.ID, userID snowflake.ID, ratingRole snowflake.ID) error {

	err = client.Rest.AddMemberRole(guildID, userID, ratingRole)
	return fmt.Errorf("error adding role: %w", err)
}

func linkRatingHandler(author *fluxer.User, message *fluxer.Message, args []string) error {
	return fmt.Errorf("unimplemented")
}

func getRating(steamID string) (int, error) {

	url := fmt.Sprintf(
		"https://api-public.cs-prod.leetify.com/v3/profile?steam64_id=%s",
		steamID,
	)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	// fmt.Println(string(body))

	var profile Profile
	err = json.Unmarshal(body, &profile)
	if err != nil {
		return 0, err
	}

	rating := profile.Ranks.Premier
	fmt.Println("Premier Rating:", rating)

	return rating, nil
}
