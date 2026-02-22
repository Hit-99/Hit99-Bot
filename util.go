package main

import (
	"strings"

	"github.com/disgoorg/snowflake/v2"
	"github.com/fluxergo/fluxergo/fluxer"
)

func getMessageById(channelID snowflake.ID, messageID snowflake.ID) (*fluxer.Message, error) {
	msg, err := client.Rest.GetMessage(channelID, messageID)
	return msg, err
}

func fixAvatarURL(url string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(url, "https://cdn.discordapp.com/", "https://fluxerusercontent.com/"),
		".png", ".webp",
	)
}
