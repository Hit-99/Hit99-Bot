package main

import (
	"context"
	"fmt"
	"os"

	"github.com/disgoorg/disgo"
	dbot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	devents "github.com/disgoorg/disgo/events"
	dgateway "github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
)

var dClient *dbot.Client

func initDiscord() *dbot.Client {
	dClient, err = disgo.New(os.Getenv("DISCORD_BOT_TOKEN"),
		dbot.WithGatewayConfigOpts(
			dgateway.WithIntents(
				dgateway.IntentGuildMessages,
				dgateway.IntentMessageContent,
			),
		),
		dbot.WithEventListenerFunc(dOnMessageCreate),
		dbot.WithEventListenerFunc(dOnReady),
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

func dOnReady(event *devents.Ready) {

	dlogid := snowflake.MustParse(os.Getenv("DLOG_CHANNEL_ID"))

	dstartmsg := discord.NewMessageCreate().WithContent("Hit-99 Bot has started!")

	_, err = dClient.Rest.CreateMessage(dlogid, dstartmsg)
	if err != nil {
		fmt.Printf("error while creating message: %s\n", err)
	}

	fmt.Println("the discord api works!!!")

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
