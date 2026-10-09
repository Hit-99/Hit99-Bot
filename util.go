package main

import (
	"encoding/json"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/fluxergo/fluxergo/fluxer"
)

func getMessageById(channelID snowflake.ID, messageID snowflake.ID) (*fluxer.Message, error) {
	msg, err := fClient.Rest.GetMessage(channelID, messageID)
	return msg, err
}

func dGetMessageById(channelID snowflake.ID, messageID snowflake.ID) (*discord.Message, error) {
	msg, err := dClient.Rest.GetMessage(channelID, messageID)
	return msg, err
}

func fixAvatarURL(url string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(url, "https://cdn.discordapp.com/", "https://fluxerusercontent.com/"),
		".png", ".webp",
	)
}

func toFloat64(n json.Number) float64 {
	f, err := n.Float64()
	if err != nil {
		return 0
	}
	return f
}

func getMatchResult(teamNumber, ctScoreNumber, tScoreNumber json.Number) (int64, int64, int, error) {
	team, err := teamNumber.Int64()
	if err != nil {
		return 0, 0, colorMatchOther, err
	}

	ctScore, err := ctScoreNumber.Int64()
	if err != nil {
		return 0, 0, colorMatchOther, err
	}

	tScore, err := tScoreNumber.Int64()
	if err != nil {
		return 0, 0, colorMatchOther, err
	}

	var playerScore, opponentScore int64

	switch team {
	case 2:
		playerScore = ctScore
		opponentScore = tScore
	case 3:
		playerScore = tScore
		opponentScore = ctScore
	default:
		return 0, 0, colorMatchOther, nil
	}

	switch {
	case playerScore > opponentScore:
		return playerScore, opponentScore, colorMatchWin, nil
	case playerScore < opponentScore:
		return playerScore, opponentScore, colorMatchLoss, nil
	default:
		return playerScore, opponentScore, colorMatchOther, nil
	}
}
