package main

import (
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
