package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/snowflake/v2"
	"github.com/fluxergo/fluxergo"
	"github.com/fluxergo/fluxergo/bot"
	"github.com/fluxergo/fluxergo/events"
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

var client *bot.Client
var err error

func main() {
	fmt.Println("Starting...")

	client, err = fluxergo.New(os.Getenv("FLUXER_BOT_TOKEN"),
		bot.WithGatewayConfigOpts(
			gateway.WithPresenceOpts(
				gateway.WithOnlineStatus(fluxer.OnlineStatusOnline)),
		),
		bot.WithEventListenerFunc(onready),
		bot.WithEventListenerFunc(commands),
		bot.WithEventListenerFunc(autorole),
		bot.WithEventListenerFunc(updateReactionRoles),
		bot.WithEventListenerFunc(userJoinEvent),
		// bot.WithEventListenerFunc(userLeaveEvent),
		bot.WithEventListenerFunc(userMsgSendEvent),
		// bot.WithEventListenerFunc(userMsgDelEvent),
		bot.WithEventListenerFunc(userMsgEditEvent),
	)

	if err != nil {
		fmt.Printf("error while building bot instance: %s\n", err)
		return
	}

	err = client.OpenGateway(context.TODO())
	if err != nil {
		fmt.Printf("error while connecting to fluxer: %s\n", err)
	}

	defer client.Close(context.TODO())

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func onready(event *events.Ready) {

	logid := snowflake.MustParse(os.Getenv("LOG_CHANNEL_ID"))

	startmsg := fluxer.NewMessageCreate().WithContent("Hit-99 Bot has started!")

	_, err = client.Rest.CreateMessage(logid, startmsg)
	if err != nil {
		fmt.Printf("error while creating message: %s\n", err)
	}

	fmt.Println("the api works!!!")

	initCommands()
	initReactionRoles()
}
