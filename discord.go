package main

import (
	"context"
	"fmt"
	"os"

	"github.com/disgoorg/disgo"
	dbot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	devents "github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

var dClient *dbot.Client

func initDiscord() *dbot.Client {
	dClient, err = disgo.New(os.Getenv("DISCORD_BOT_TOKEN"),
		dbot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds|
					gateway.IntentGuildMessages|
					gateway.IntentMessageContent|
					gateway.IntentGuildMessageReactions|
					gateway.IntentGuildMembers,
			),
			gateway.WithPresenceOpts(
				gateway.WithOnlineStatus(discord.OnlineStatusOnline),
			),
		),
		dbot.WithEventListenerFunc(dOnMessageCreate),
		dbot.WithEventListenerFunc(dOnReady),
		dbot.WithEventListenerFunc(dcommands),
		dbot.WithEventListenerFunc(dAutorole),
		dbot.WithEventListenerFunc(dUpdateReactionRoles),
		dbot.WithEventListenerFunc(dUserJoinEvent),
		// bot.WithEventListenerFunc(dUserLeaveEvent),
		dbot.WithEventListenerFunc(dUserMsgSendEvent),
		// bot.WithEventListenerFunc(dUserMsgDelEvent),
		dbot.WithEventListenerFunc(dUserMsgEditEvent),
	)

	if err != nil {
		fmt.Printf("error while building discord bot instance: %s\n", err)
		return nil
	}

	err = dClient.OpenGateway(context.TODO())
	if err != nil {
		fmt.Printf("error while connecting to discord: %s\n", err)
	}

	return dClient

}

func dOnMessageCreate(event *devents.MessageCreate) {
	if event.Message.Author.Bot {
		return
	}

	attachements := event.Message.Attachments

	channelID := event.Message.ChannelID.String()
	content := event.Message.Content
	username := event.Message.Author.Username
	avatarUrl := *event.Message.Author.AvatarURL()

	err = sendFluxerWebhookBridgeMessage(channelID, content, username, avatarUrl, attachements)
	if err != nil {
		fmt.Printf("error while sending message to fluxer: %s\n", err)
	}
}
