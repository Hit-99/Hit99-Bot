package main

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

var steamID string
var matchChannelID snowflake.ID
var dMatchChannelID snowflake.ID

// checks for new matches from a steam user list and sends them in the specified channel
func csMatchListener(playerstats LeetifyProfile) {
	matchChannelID = snowflake.MustParse("1475607334483988934")
	dMatchChannelID = snowflake.MustParse("1499780625966829598")
	for _, match := range playerstats.RecentMatches {
		matchID := match.ID
		recentMatch := playerstats.RecentMatches[0].ID
		steamID = playerstats.SteamID

		time.Sleep(1 * time.Second)

		ifMatchExist, _ := ifMatchExistsForUser(steamID, matchID)
		if !ifMatchExist {
			matchesCreateEntry(steamID, recentMatch)
			matchStatsEmbed(matchChannelID, playerstats, matchID)
			dMatchStatsEmbed(dMatchChannelID, playerstats, matchID)
			return
		} else if matchID == recentMatch {
			return
		} else { // will send all previous in reverse chronological order
			matchesCreateEntry(steamID, matchID)
			matchStatsEmbed(matchChannelID, playerstats, matchID)
			dMatchStatsEmbed(dMatchChannelID, playerstats, matchID)
			if matchID == recentMatch {
				return
			}
		}
	}
}

// if bot is reconnecting, it will still add to db even if message is never sent
