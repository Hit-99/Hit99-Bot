package main

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/fluxergo/fluxergo/fluxer"
	"github.com/fluxergo/fluxergo/rest"
)

var fluxerDiscordBridge = map[string]snowflake.ID{
	"708918939710652487": snowflake.ID(1474820699636891678),
}

var fluxerWebhookTokens = map[snowflake.ID]string{
	snowflake.ID(1474820699636891678): "GqNsEkoR4ia3Vv64oTbk93dp5sV1qL4QS00XcCDySMxYXMuwLu48X99muw5KU41a",
}

func webhooktesthandler(caller *fluxer.User, message *fluxer.Message, args []string) error {
	if caller.ID != snowflake.MustParse("1473029871908167709") {
		return nil
	}

	// webhookcreate := fluxer.WebhookCreate{
	// 	Name: "test",
	// }

	// webhook, err := client.Rest.CreateWebhook(message.ChannelID, webhookcreate)

	// fmt.Printf("ID: %v\nName: %s\nAvatar: %v\nChannelID: %v\nGuildID: %v\nToken: %s\nApplicationID: %v\nUser: %v\n",
	// 	webhook.ID().String(), webhook.Name(), webhook.Avatar, webhook.ChannelID, webhook.GuildID, webhook.Token, webhook.ApplicationID, webhook.User)

	webhookmsg := fluxer.WebhookMessageCreate{
		Content:   args[1],
		Username:  caller.Username,
		AvatarURL: fixAvatarURL(*caller.AvatarURL()),
	}

	_, err = fClient.Rest.CreateWebhookMessage(snowflake.ID(1474820699636891678), "GqNsEkoR4ia3Vv64oTbk93dp5sV1qL4QS00XcCDySMxYXMuwLu48X99muw5KU41a", webhookmsg, rest.CreateWebhookMessageParams{})

	return err
}

func sendFluxerWebhookBridgeMessage(discordChannelID string, discordMsgContent string, discordUsername string, discordAvatarURL string, attachements []discord.Attachment) error {

	webhookmsg := fluxer.WebhookMessageCreate{
		Content:   discordMsgContent,
		Username:  discordUsername,
		AvatarURL: discordAvatarURL,
	}

	if len(attachements) > 0 {

		var fluxerEmbeds []fluxer.Embed
		for _, attachment := range attachements {
			attachmentURL := attachment.URL
			fmt.Println(attachmentURL)
			fluxerEmbeds = append(fluxerEmbeds, fluxer.Embed{
				Title: "Media",
				Image: &fluxer.EmbedResource{
					URL: attachmentURL,
				},
			})
		}
		webhookmsg.Embeds = fluxerEmbeds

	}

	_, err = fClient.Rest.CreateWebhookMessage(fluxerDiscordBridge[discordChannelID], fluxerWebhookTokens[fluxerDiscordBridge[discordChannelID]], webhookmsg, rest.CreateWebhookMessageParams{})
	return err
}
