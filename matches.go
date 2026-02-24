package main

import (
	"github.com/disgoorg/snowflake/v2"
	"github.com/fluxergo/fluxergo/events"
)

var steamID string // make array/json file
var matchChannelID snowflake.ID

// hypersteamID = "76561198418069261"
// nevanaSteamID = "76561198881070768"

// checks for new matches from a steam user list and sends them in the specified channel

func csMatchListener(event *events.Ready) {
	steamID = "76561198418069261"
	matchChannelID = snowflake.MustParse("1474077652385337399") // staff bot commands temporarily
	playerstats, _ := getLeetifyStats(steamID)
	matchStatsEmbed(matchChannelID, playerstats)
	// need a check to see if match data has already been sent
}

// need a list of steam IDs
// check every so often (5 mins?) if a user has completed a match recently (within the check timeframe)
// if there is a match, send the match details in the specified channel

// need to get current date and time, compare match date and time to current time, and if its within the check timeframe, it will post the stats to the channel

// you link fluxer to steam

// db linking fluxer to steam

// each loop check for new steam ids

// table 1 -> fluxer id, steam id
// table 2 -> match id, steam id

// loop will take steam ids, send api request, and check match ids against table 2. if match id is not in table, it adds it and sends a message
