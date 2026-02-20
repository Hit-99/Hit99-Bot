package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
	rating := getRating(args[1])

	ratingmsg := fluxer.NewMessageCreate().WithContent(fmt.Sprintf("Your premier rating is %d", rating))

	_, err = client.Rest.CreateMessage(message.ChannelID, ratingmsg)

	return err

}

func linkRatingHandler(author *fluxer.User, message *fluxer.Message, args []string) error {
	return fmt.Errorf("unimplemented")
}

func getRating(steamID string) int {

	url := fmt.Sprintf(
		"https://api-public.cs-prod.leetify.com/v3/profile?steam64_id=%s",
		steamID,
	)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// fmt.Println(string(body))

	var profile Profile
	err = json.Unmarshal(body, &profile)
	if err != nil {
		panic(err)
	}

	rating := profile.Ranks.Premier
	fmt.Println("Premier Rating:", rating)

	return rating
}
