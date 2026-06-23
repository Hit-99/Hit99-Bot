package main

import (
	"fmt"

	devents "github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/fluxergo/fluxergo/events"
)

// give users member role on join
func autorole(event *events.GuildMemberJoin) {
	fmt.Println("user joined")
	memberRoleID := snowflake.MustParse("1474078482417119356")
	err = fClient.Rest.AddMemberRole(event.GuildID, event.Member.User.ID, memberRoleID)
}

func dAutorole(event *devents.GuildMemberJoin) {
	fmt.Println("user joined")
	memberRoleID := snowflake.MustParse("712775854878359563")
	err = dClient.Rest.AddMemberRole(event.GuildID, event.Member.User.ID, memberRoleID)
}
