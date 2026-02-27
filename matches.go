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
	fmt.Println("(debug) matchID: ", apiMatchID) // debug
	steamID = playerstats.SteamID

	ifMatchExist, _ := ifMatchExistsForUser(steamID, apiMatchID)
	if !ifMatchExist {
		matchesCreateEntry(steamID, apiMatchID)
		matchStatsEmbed(matchChannelID, playerstats)
	}
	// fmt.Println("(debug) matches updated")
}
