package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/fluxergo/fluxergo"
	"github.com/fluxergo/fluxergo/bot"
	"github.com/fluxergo/fluxergo/fluxer"
	"github.com/fluxergo/fluxergo/gateway"
	"github.com/joho/godotenv"
)

func init() {
	fmt.Println("Init...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

var fClient *bot.Client

func initFluxer() *bot.Client {
	fClient, err = fluxergo.New(os.Getenv("FLUXER_BOT_TOKEN"),
		bot.WithGatewayConfigOpts(
			gateway.WithPresenceOpts(
				gateway.WithOnlineStatus(fluxer.OnlineStatusOnline)),
		),
		bot.WithEventListenerFunc(fOnReady),
		bot.WithEventListenerFunc(commands),
		bot.WithEventListenerFunc(autorole),
		bot.WithEventListenerFunc(updateReactionRoles),
		bot.WithEventListenerFunc(userJoinEvent),
		bot.WithEventListenerFunc(userLeaveEvent),
		bot.WithEventListenerFunc(userMsgSendEvent),
		bot.WithEventListenerFunc(userMsgEditEvent),
		bot.WithEventListenerFunc(userMsgDelEvent),
		// bot.WithEventListenerFunc(userMsgReplyEvent),
		bot.WithEventListenerFunc(roleCreateEvent),
		bot.WithEventListenerFunc(roleUpdateEvent),
		bot.WithEventListenerFunc(roleDeleteEvent),
		bot.WithEventListenerFunc(channelCreateEvent),
		bot.WithEventListenerFunc(channelUpdateEvent),
		bot.WithEventListenerFunc(channelDeleteEvent),
	)
	if err != nil {
		fmt.Printf("error while building bot instance: %s\n", err)
		return nil
	}

	err = fClient.OpenGateway(context.TODO())
	if err != nil {
		fmt.Printf("error while connecting to fluxer: %s\n", err)
	}
	return fClient
}
