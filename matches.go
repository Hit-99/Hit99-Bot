package main

import (
	"fmt"

	"github.com/disgoorg/snowflake/v2"
)

var steamID string
var matchChannelID snowflake.ID

// checks for new matches from a steam user list and sends them in the specified channel

func csMatchListener(playerstats LeetifyProfile) {
	matchChannelID = snowflake.MustParse("1475607334483988934")
	apiMatchID := playerstats.RecentMatches[0].ID
	steamID = playerstats.SteamID
	if matchesContainsMatchID(apiMatchID) == true {
		if matchesContainsSteamID(steamID) == true {
			// if match id and steam id in db, do nothing
			// fmt.Println("Match already in db")
			return
		}
		matchesCreateEntry(steamID, apiMatchID) // if match id in db but steam id not, create new entry
		matchStatsEmbed(matchChannelID, playerstats)
		// fmt.Println("Adding match to db")
	} else {
		matchesCreateEntry(steamID, apiMatchID) // if match id and steam id not in db, create new entry
		matchStatsEmbed(matchChannelID, playerstats)
		// fmt.Println("Adding match to db")
	}
	fmt.Println("matches updated")
}
