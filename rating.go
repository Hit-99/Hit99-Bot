package main

import "github.com/disgoorg/snowflake/v2"

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
