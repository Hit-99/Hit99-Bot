package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type VanityResponse struct {
	Response struct {
		SteamID string `json:"steamid"`
		Success int    `json:"success"`
		Message string `json:"message"`
	} `json:"response"`
}

func resolveVanity(apiKey string, name string) (VanityResponse, error) {
	fmt.Println("(debug) resolveVanity name: ", name)
	url := fmt.Sprintf(
		"https://api.steampowered.com/ISteamUser/ResolveVanityURL/v1/?key=%s&vanityurl=%s",
		apiKey,
		name,
	)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("error while getting response from steam: %s\n", err)
		return VanityResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error while reading from steam: %s\n", err)
		return VanityResponse{}, err
	}

	var data VanityResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Printf("error while unmarshalling json from steam: %s\n", err)
		return VanityResponse{}, err
	}

	if data.Response.Success != 1 {
		return VanityResponse{}, fmt.Errorf("could not resolve vanity URL")
	}

	return data, nil
}

func getSteamID(name string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("STEAM_API_KEY")
	if apiKey == "" {
		log.Fatal("STEAM_API_KEY not set")
	}

	id, err := resolveVanity(apiKey, name)
	if err != nil {
		fmt.Println("Error:", err)
	}
	steamID = id.Response.SteamID
	fmt.Println("(debug) SteamID64:", steamID)
	return steamID, err
}
