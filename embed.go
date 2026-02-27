package main

import (
	"fmt"
	"strconv"
	"time"

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
	_, err = client.Rest.CreateMessage(channelID, embedOut)
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
	_, err = client.Rest.CreateMessage(channelID, embedOut)
	return err
}

func matchStatsEmbed(channelID snowflake.ID, matchStats LeetifyProfile) error {
	switch matchStats.RecentMatches[0].Outcome {
	case "loss":
		embedColor = colorMatchLoss
	case "win":
		embedColor = colorMatchWin
	default:
		embedColor = colorMatchOther
	}
	finishedAtStr := matchStats.RecentMatches[0].FinshedAt
	parsedTime, err := time.Parse(time.RFC3339, finishedAtStr)
	if err != nil {
		return err
	}

	leetifyReformat, _ := matchStats.RecentMatches[0].LeetifyRating.Float64()
	playerRank, _ := matchStats.RecentMatches[0].Rank.Int64()
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
		Title: "View on Leetify",
		Color: embedColor, // make if statement to swap win and loss colors
		Author: &fluxer.EmbedAuthor{
			Name: fmt.Sprintf("%s  %s  ->  %s", matchStats.RecentMatches[0].Score, matchStats.RecentMatches[0].MapName, matchStats.UserName),
		},
		Description: fmt.Sprintf("**Rank:** %s\n**Leetify Rating:** %v\n**Preaim:** %s\n**Reaction Time (MS):** %s\n**Accuracy Enemy Spotted:** %s\n**Accuracy Head:** %s\n**Spray Accuracy:** %s",
			playerRankName,
			(leetifyReformat * 100),
			matchStats.RecentMatches[0].Preaim,
			matchStats.RecentMatches[0].ReactionTimeMS,
			matchStats.RecentMatches[0].AccuracyEnemySpotted,
			matchStats.RecentMatches[0].AccuracyHead,
			matchStats.RecentMatches[0].SprayAccuracy),
		Timestamp: &parsedTime,
		URL:       fmt.Sprintf("https://leetify.com/app/match-details/%s/your-match", matchStats.RecentMatches[0].ID),
		Thumbnail: &fluxer.EmbedResource{
			URL:    "https://cloud.hy7.dev/apps/files_sharing/publicpreview/dooMgMXNSf3Q4r2?file=/&fileId=4110&x=1920&y=1080&a=true&etag=535ca46af95fa6a0cca81a465677911d",
			Height: 115,
			Width:  270,
		},
		Footer: &fluxer.EmbedFooter{
			Text:    matchStats.RecentMatches[0].ID,
			IconURL: "https://fluxerusercontent.com/attachments/1473793058206990390/1475614229764805051/Artboard_1.png",
		},
	}
	embedOut := fluxer.NewMessageCreate().WithContent("").WithEmbeds(embed)
	_, err = client.Rest.CreateMessage(channelID, embedOut)
	return err
}
