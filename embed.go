package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/disgoorg/disgo/discord"
	devents "github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/fluxergo/fluxergo/events"
	"github.com/fluxergo/fluxergo/fluxer"
)

const (
	colorHit99      = 0x4caf50
	colorError      = 0xf44336
	colorEdit       = 0xff9800
	colorMatchWin   = 0x4caf50
	colorMatchLoss  = 0xf44336
	colorMatchOther = 0x808080
)

var embedColor int

func delMsgEmbed(channelID snowflake.ID, embedMsg *events.GuildMessageDelete) error {
	embed := fluxer.Embed{
		Author: &fluxer.EmbedAuthor{
			Name:    embedMsg.Message.Author.Username,
			IconURL: *embedMsg.Message.Author.AvatarURL(),
		},
		Color:       0xf44336,
		Description: embedMsg.Message.Content, // doesnt like to be empty
		Timestamp:   &embedMsg.Message.CreatedAt,
	}
	embedOut := fluxer.NewMessageCreate().WithContent("").WithEmbeds(embed)
	_, err = fClient.Rest.CreateMessage(channelID, embedOut)
	return err
}

func editMsgEmbed(channelID snowflake.ID, embedMsg *events.GuildMessageUpdate) error {
	embed := fluxer.Embed{
		Author: &fluxer.EmbedAuthor{
			Name:    embedMsg.Message.Author.Username,
			IconURL: *embedMsg.Message.Author.AvatarURL(),
		},
		Color:       0xf44336,
		Description: "Before: " + embedMsg.OldMessage.Content + "After: " + embedMsg.Message.Content,
		Timestamp:   &embedMsg.Message.CreatedAt,
	}
	embedOut := fluxer.NewMessageCreate().WithContent("").WithEmbeds(embed)
	_, err = fClient.Rest.CreateMessage(channelID, embedOut)
	return err
}

func dEditMsgEmbed(channelID snowflake.ID, embedMsg *devents.GuildMessageUpdate) error {
	embed := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name:    embedMsg.Message.Author.Username,
			IconURL: *embedMsg.Message.Author.AvatarURL(),
		},
		Color:       0xf44336,
		Description: "Before: " + embedMsg.OldMessage.Content + "After: " + embedMsg.Message.Content,
		Timestamp:   &embedMsg.Message.CreatedAt,
	}
	embedOut := discord.NewMessageCreate().WithContent("").WithEmbeds(embed)
	_, err = dClient.Rest.CreateMessage(channelID, embedOut)
	return err
}

func matchStatsEmbed(channelID snowflake.ID, matchStats CSMetricsProfile, matchID string) error {
	var recentMatch int
	for i := range matchStats.Matches {
		if matchStats.Matches[i].MatchUUID == matchID {
			recentMatch = i
		}
	}

	match := matchStats.Matches[recentMatch]

	playerScore, opponentScore, resultColor, err := getMatchResult(
		match.Team,
		match.CTScore,
		match.TScore,
	)
	if err != nil {
		return err
	}

	team, err := match.Team.Int64()
	if err != nil {
		return err
	}

	ctScore, err := match.CTScore.Int64()
	if err != nil {
		return err
	}

	tScore, err := match.TScore.Int64()
	if err != nil {
		return err
	}

	switch {
	case ctScore == tScore:
		embedColor = colorMatchOther

	case (team == 2 && tScore > ctScore) ||
		(team == 3 && ctScore > tScore):
		embedColor = colorMatchWin

	case team == 2 || team == 3:
		embedColor = colorMatchLoss

	default:
		embedColor = colorMatchOther
	}

	finishedAtStr := matchStats.Matches[recentMatch].MatchDate
	parsedTime, err := time.Parse(time.RFC3339, finishedAtStr)
	if err != nil {
		return err
	}

	playerRank, _ := matchStats.Matches[recentMatch].Rank.Int64()
	var playerRankName string
	var playerRankMap = map[string]string{
		"0":  "Unranked",
		"1":  "Silver 1",
		"2":  "Silver II",
		"3":  "Silver III",
		"4":  "Silver IV",
		"5":  "Silver Elite",
		"6":  "Silver Elite Master",
		"7":  "Gold Nova I",
		"8":  "Gold Nova II",
		"9":  "Gold Nova III",
		"10": "Gold Nova Master",
		"11": "Master Guardian I",
		"12": "Master Guardian II",
		"13": "Master Guardian Elite",
		"14": "Destinguished Master Guardian",
		"15": "Legendary Eagle",
		"16": "Legendary Eagle Master",
		"17": "Supreme Master First Class",
		"18": "The Global Elite",
	}

	if playerRank < 1000 {
		playerRankName = playerRankMap[strconv.Itoa(int(playerRank))]
	} else {
		playerRankName = strconv.Itoa(int(playerRank))

	}

	embed := fluxer.Embed{
		Title: "View on CSMetrics",
		Color: resultColor,
		Author: &fluxer.EmbedAuthor{
			Name: fmt.Sprintf("[%d - %d]  %s  ->  %s",
				playerScore,
				opponentScore,
				matchStats.Matches[recentMatch].Map,
				matchStats.UserName),
		},
		Description: fmt.Sprintf("**Rank:** %s\n**Kills:** %s **Deaths:** %s **Assists:** %s\n**KAST:** %.1f%%\n**Crosshair Placement:** %.1f°\n**Time To Damage:** %.0fms",
			playerRankName,
			matchStats.Matches[recentMatch].Kills,
			matchStats.Matches[recentMatch].Deaths,
			matchStats.Matches[recentMatch].Assists,
			toFloat64(matchStats.Matches[recentMatch].KASTRounds)/toFloat64(matchStats.Matches[recentMatch].RoundsPlayed)*100,
			toFloat64(matchStats.Matches[recentMatch].MedianCrosshair),
			toFloat64(matchStats.Matches[recentMatch].AverageTimeToDamage)),
		Timestamp: &parsedTime,
		URL:       fmt.Sprintf("https://csmetrics.app/match/%s", matchStats.Matches[recentMatch].MatchUUID),
		Thumbnail: &fluxer.EmbedResource{
			URL:    "https://cloud.hy7.dev/apps/files_sharing/publicpreview/dooMgMXNSf3Q4r2?file=/&fileId=4110&x=1920&y=1080&a=true&etag=535ca46af95fa6a0cca81a465677911d",
			Height: 115,
			Width:  270,
		},
		Footer: &fluxer.EmbedFooter{
			Text:    matchStats.Matches[recentMatch].MatchUUID,
			IconURL: "https://fluxerusercontent.com/attachments/1473793058206990390/1475614229764805051/Artboard_1.png",
		},
	}
	embedOut := fluxer.NewMessageCreate().WithContent("").WithEmbeds(embed)
	_, err = fClient.Rest.CreateMessage(channelID, embedOut)
	return err
}

func dMatchStatsEmbed(channelID snowflake.ID, matchStats CSMetricsProfile, matchID string) error {
	var recentMatch int
	for i := range matchStats.Matches {
		if matchStats.Matches[i].MatchUUID == matchID {
			recentMatch = i
		}
	}

	match := matchStats.Matches[recentMatch]

	playerScore, opponentScore, resultColor, err := getMatchResult(
		match.Team,
		match.CTScore,
		match.TScore,
	)
	if err != nil {
		return err
	}

	team, err := match.Team.Int64()
	if err != nil {
		return err
	}

	ctScore, err := match.CTScore.Int64()
	if err != nil {
		return err
	}

	tScore, err := match.TScore.Int64()
	if err != nil {
		return err
	}

	switch {
	case ctScore == tScore:
		embedColor = colorMatchOther

	case (team == 2 && tScore > ctScore) ||
		(team == 3 && ctScore > tScore):
		embedColor = colorMatchWin

	case team == 2 || team == 3:
		embedColor = colorMatchLoss

	default:
		embedColor = colorMatchOther
	}

	finishedAtStr := matchStats.Matches[recentMatch].MatchDate
	parsedTime, err := time.Parse(time.RFC3339, finishedAtStr)
	if err != nil {
		return err
	}

	playerRank, _ := matchStats.Matches[recentMatch].Rank.Int64()
	var playerRankName string
	var playerRankMap = map[string]string{
		"0":  "Unranked",
		"1":  "Silver 1",
		"2":  "Silver II",
		"3":  "Silver III",
		"4":  "Silver IV",
		"5":  "Silver Elite",
		"6":  "Silver Elite Master",
		"7":  "Gold Nova I",
		"8":  "Gold Nova II",
		"9":  "Gold Nova III",
		"10": "Gold Nova Master",
		"11": "Master Guardian I",
		"12": "Master Guardian II",
		"13": "Master Guardian Elite",
		"14": "Destinguished Master Guardian",
		"15": "Legendary Eagle",
		"16": "Legendary Eagle Master",
		"17": "Supreme Master First Class",
		"18": "The Global Elite",
	}

	if playerRank < 1000 {
		playerRankName = playerRankMap[strconv.Itoa(int(playerRank))]
	} else {
		playerRankName = strconv.Itoa(int(playerRank))

	}

	embed := discord.Embed{
		Title: "View on CSMetrics",
		Color: resultColor,
		Author: &discord.EmbedAuthor{
			Name: fmt.Sprintf("[%d - %d]  %s  ->  %s",
				playerScore,
				opponentScore,
				match.Map,
				matchStats.UserName,
			),
		},
		Description: fmt.Sprintf("**Rank:** %s\n**Kills:** %s **Deaths:** %s **Assists:** %s\n**KAST:** %.1f%%\n**Crosshair Placement:** %.1f°\n**Time To Damage:** %.0fms",
			playerRankName,
			matchStats.Matches[recentMatch].Kills,
			matchStats.Matches[recentMatch].Deaths,
			matchStats.Matches[recentMatch].Assists,
			toFloat64(matchStats.Matches[recentMatch].KASTRounds)/toFloat64(matchStats.Matches[recentMatch].RoundsPlayed)*100,
			toFloat64(matchStats.Matches[recentMatch].MedianCrosshair),
			toFloat64(matchStats.Matches[recentMatch].AverageTimeToDamage)),
		Timestamp: &parsedTime,
		URL:       fmt.Sprintf("https://csmetrics.app/match/%s", matchStats.Matches[recentMatch].MatchUUID),
		Thumbnail: &discord.EmbedResource{
			URL:    "https://cloud.hy7.dev/apps/files_sharing/publicpreview/dooMgMXNSf3Q4r2?file=/&fileId=4110&x=1920&y=1080&a=true&etag=535ca46af95fa6a0cca81a465677911d",
			Height: 115,
			Width:  270,
		},
		Footer: &discord.EmbedFooter{
			Text:    matchStats.Matches[recentMatch].MatchUUID,
			IconURL: "https://fluxerusercontent.com/attachments/1473793058206990390/1475614229764805051/Artboard_1.png",
		},
	}
	embedOut := discord.NewMessageCreate().WithContent("").WithEmbeds(embed)
	_, err = dClient.Rest.CreateMessage(channelID, embedOut)
	return err
}
