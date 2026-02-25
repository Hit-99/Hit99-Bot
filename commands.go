package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/disgoorg/snowflake/v2"
	"github.com/fluxergo/fluxergo/events"
	"github.com/fluxergo/fluxergo/fluxer"
)

var (
	prefix        string
	prefixPattern *regexp.Regexp
)

type Command struct {
	Name    string
	Handler func(caller *fluxer.User, message *fluxer.Message, args []string) error
}

var Commands = []Command{
	{
		Name:    "ping",
		Handler: pingpong, //pong
	},
	{
		Name:    "rating",
		Handler: getPremierRatingHandler, //get premier rating, return to user
	},
	{
		Name:    "link",
		Handler: linkSteamHandler, // link steam and fluxer ids
	},
	{
		Name:    "getroles",
		Handler: getallroleshandler,
	},
	{
		Name:    "webhooktest",
		Handler: webhooktesthandler,
	},
	{
		Name:    "stats",
		Handler: getStatsHandler, //get basic stats, return to user
	},
}

func getallroleshandler(caller *fluxer.User, message *fluxer.Message, args []string) error {

	roles, _ := client.Rest.GetRoles(*message.GuildID)
	for _, role := range roles {
		fmt.Printf("Role: %s, ID: %s\n", role.Name, role.ID)

	}
	return nil
}

func pingpong(caller *fluxer.User, message *fluxer.Message, args []string) error {

	pongmsg := fluxer.NewMessageCreate().WithContent("pong")
	_, err = client.Rest.CreateMessage(message.ChannelID, pongmsg)

	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}

func initCommands() {
	prefix = os.Getenv("PREFIX")
	if prefix == "" {
		prefix = "!"
	}

	prefixPattern = regexp.MustCompile(`(?i)^\s*` + regexp.QuoteMeta(prefix) + `\s*`)
}

func findCommmand(name string) *Command {
	for _, command := range Commands {
		if name == command.Name {
			return &command
		}
	}
	return nil
}

func commands(event *events.MessageCreate) {

	message := &event.Message

	if message.Author.ID == client.ID() {
		return
	}

	// !getRating 23409823490832
	content := message.Content

	match := prefixPattern.FindString(content)

	if match != "" {
		// !
		args := strings.Fields(content[len(match):])
		// [getRating, 23409823490832]
		if len(args) == 0 {
			return
		}

		commandName := strings.ToLower(args[0])
		// getRating

		command := findCommmand(commandName)
		// 	{
		// 	Name:    "getRating",
		// 	Handler: getRatingHandler, //get rating, return to user
		//  },

		if command == nil {

			cmdnotfounderrmsg := fluxer.NewMessageCreate().WithContent("cmd not found")

			_, err = client.Rest.CreateMessage(message.ChannelID, cmdnotfounderrmsg)
			if err != nil {
				fmt.Printf("error while finding command: %s\n", err)
			}
			return
		}

		author, err := getUserByID(message.Author.ID)
		if err != nil {
			return
		}

		err = command.Handler(author, message, args)
		if err != nil {
			fmt.Printf("error while executing command: %s\n", err)
		}

	}

}

func getUserByID(userID snowflake.ID) (*fluxer.User, error) {
	user, err := client.Rest.GetUser(userID)
	if err != nil {
		fmt.Printf("error while fetching user: %s\n", err)
		return nil, err
	}
	return user, nil
}
