package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo/discord"
	devents "github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/fluxergo/fluxergo/events"
	"github.com/fluxergo/fluxergo/fluxer"
)

var err error

func main() {
	fmt.Println("Starting...")

	fClient := initFluxer()
	dClient := initDiscord()

	fmt.Println("initialized discord bridge") // debug

	defer fClient.Close(context.TODO())
	defer dClient.Close(context.TODO())

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func fOnReady(event *events.Ready) {

	logid := snowflake.MustParse(os.Getenv("LOG_CHANNEL_ID"))

	startmsg := fluxer.NewMessageCreate().WithContent("Hit-99 Bot has started!")

	_, err = fClient.Rest.CreateMessage(logid, startmsg)
	if err != nil {
		fmt.Printf("error while creating message: %s\n", err)
	}

	fmt.Println("the fluxer api works!!!")

	initCommands()
	fmt.Println("initialized fluxer commands")
	initReactionRoles()
	fmt.Println("initialized fluxer reaction roles")
	go initLeetifyStatsLoop()
	fmt.Println("initialized looped functions")

}

func dOnReady(event *devents.Ready) {

	dlogid := snowflake.MustParse(os.Getenv("DLOG_CHANNEL_ID"))

	dstartmsg := discord.NewMessageCreate().WithContent("Hit-99 Bot has started!")

	_, err = dClient.Rest.CreateMessage(dlogid, dstartmsg)
	if err != nil {
		fmt.Printf("error while creating message: %s\n", err)
	}

	fmt.Println("the discord api works!!!")

	dInitCommands()
	fmt.Println("initialized discord commands")
	dInitReactionRoles()
	fmt.Println("initialized discord reaction roles")
}
