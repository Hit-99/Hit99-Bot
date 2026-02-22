package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	dbot "github.com/disgoorg/disgo/bot"
	devents "github.com/disgoorg/disgo/events"
	dgateway "github.com/disgoorg/disgo/gateway"
)

func initDiscord() {
	dClient, derr := disgo.New(os.Getenv("DISCORD_BOT_TOKEN"),
		dbot.WithGatewayConfigOpts(
			dgateway.WithIntents(
				dgateway.IntentGuildMessages,
				dgateway.IntentMessageContent,
			),
		),
		dbot.WithEventListenerFunc(dOnMessageCreate),
		dbot.WithEventListenerFunc(dOnReady),
	)

	if derr != nil {
		fmt.Printf("error while building discord bot instance: %s\n", derr)
		return
	}

	derr = dClient.OpenGateway(context.TODO())
	if derr != nil {
		fmt.Printf("error while connecting to discord: %s\n", derr)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func dOnReady(event *devents.Ready) {

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
