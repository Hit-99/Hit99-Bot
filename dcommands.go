package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/disgoorg/disgo/discord"
	devents "github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

var (
	dPrefix        string
	dPrefixPattern *regexp.Regexp
)

type dCommand struct {
	Name    string
	Handler func(caller *discord.User, message *discord.Message, args []string) error
}

var dCommands = []dCommand{
	{
		Name:    "ping",
		Handler: dPingPong, //pong
	},
	{
		Name:    "rating",
		Handler: dGetPremierRatingHandler, //get premier rating, return to user
	},
	{
		Name:    "link",
		Handler: dLinkSteamHandler, // link steam and discord ids
	},
	{
		Name:    "getroles",
		Handler: dGetAllRolesHandler,
	},
	{
		Name:    "stats",
		Handler: dGetStatsHandler, //get basic stats, return to user
	},
}

func dGetAllRolesHandler(caller *discord.User, message *discord.Message, args []string) error {

	roles, _ := dClient.Rest.GetRoles(*message.GuildID)
	for _, role := range roles {
		fmt.Printf("Role: %s, ID: %s\n", role.Name, role.ID)

	}
	return nil
}

func dPingPong(caller *discord.User, message *discord.Message, args []string) error {

	pongmsg := discord.NewMessageCreate().WithContent("pong")
	_, err = dClient.Rest.CreateMessage(message.ChannelID, pongmsg)

	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}

func dInitCommands() {
	prefix = os.Getenv("PREFIX")
	if prefix == "" {
		prefix = "!"
	}

	prefixPattern = regexp.MustCompile(`(?i)^\s*` + regexp.QuoteMeta(prefix) + `\s*`)
}

func dFindCommmand(name string) *dCommand {
	for _, command := range dCommands {
		if name == command.Name {
			return &command
		}
	}
	return nil
}

func dcommands(event *devents.MessageCreate) {
	message := &event.Message

	if message.Author.ID == dClient.ID() {
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

		command := dFindCommmand(commandName)
		// 	{
		// 	Name:    "getRating",
		// 	Handler: getRatingHandler, //get rating, return to user
		//  },

		if command == nil {

			cmdnotfounderrmsg := discord.NewMessageCreate().WithContent("cmd not found")

			_, err = dClient.Rest.CreateMessage(message.ChannelID, cmdnotfounderrmsg)
			if err != nil {
				fmt.Printf("error while finding command: %s\n", err)
			}
			return
		}

		author, err := dGetUserByID(message.Author.ID)
		if err != nil {
			return
		}

		err = command.Handler(author, message, args)
		if err != nil {
			fmt.Printf("error while executing command: %s\n", err)
		}

	}

}

func dGetUserByID(userID snowflake.ID) (*discord.User, error) {
	user, err := dClient.Rest.GetUser(userID)
	// throwing errors when users run commands
	if err != nil {
		fmt.Printf("error while fetching user: %s\n", err)
		return nil, err
	}
	return user, nil
}
