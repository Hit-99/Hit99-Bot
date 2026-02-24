package main

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

var steamID string
var matchChannelID snowflake.ID

// checks for new matches from a steam user list and sends them in the specified channel

func csMatchListener() {
	steamIDList := linkGetAllSteamIds()
	for _, steamID := range steamIDList {
		matchChannelID = snowflake.MustParse("1475607334483988934")
		playerstats, _ := getLeetifyStats(steamID)
		apiMatchID := playerstats.RecentMatches[0].ID
		if matchesContainsMatchID(apiMatchID) == true {
			if matchesContainsSteamID(steamID) == true {
				// if match id and steam id in db, do nothing
				return
			}
			matchesCreateEntry(steamID, apiMatchID) // if match id in db but steam id not, create new entry
			matchStatsEmbed(matchChannelID, playerstats)
		} else {
			matchesCreateEntry(steamID, apiMatchID) // if match id and steam id not in db, create new entry
			matchStatsEmbed(matchChannelID, playerstats)
		}
	}
}

// starts with bot and runs function every 5 mins

func initMatchListenter() {
	csMatchListener()
	for _ = range time.Tick(time.Minute * 5) {
		csMatchListener()
	}
	// maybe save all stats to db everytime this runs to avoid api requests?
}
