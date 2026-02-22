package main

import (
	"github.com/disgoorg/snowflake/v2"
	"github.com/fluxergo/fluxergo/events"
	"github.com/fluxergo/fluxergo/fluxer"
)

const (
	colorHit99 = 0x4caf50
	colorError = 0xf44336
	colorEdit  = 0xff9800
)

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
		Description: "After: " + embedMsg.Message.Content,
		Timestamp:   &embedMsg.Message.CreatedAt,
	}
	embedOut := fluxer.NewMessageCreate().WithContent("").WithEmbeds(embed)
	_, err = client.Rest.CreateMessage(channelID, embedOut)
	return err
}
