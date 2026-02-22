package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Profile struct {
	Ranks struct {
		Premier int `json:"premier"`
	} `json:"ranks"`
}

// gets leetify stats (need to add api key to avoid rate limiting later)

func getStats(steamID string) (Profile, error) {

	url := fmt.Sprintf(
		"https://api-public.cs-prod.leetify.com/v3/profile?steam64_id=%s",
		steamID,
	)

	resp, err := http.Get(url)
	if err != nil {
		return Profile{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Profile{}, err
	}

	var profile Profile
	err = json.Unmarshal(body, &profile)
	if err != nil {
		return Profile{}, err
	}

	return profile, nil
}

// gets premier rating for rating roles

func getRating(steamID string) (int, error) {

	profile, err := getStats(steamID)
	if err != nil {
		return 0, err
	}

	rating := profile.Ranks.Premier
	fmt.Println("Premier Rating:", rating)

	return rating, nil
}
