package main

import (
	"fmt"

	"github.com/disgoorg/snowflake/v2"
	"github.com/fluxergo/fluxergo/events"
)

// give users member role on join

func autorole(event *events.GuildMemberJoin) {
	fmt.Println("user joined")
	memberRoleID := snowflake.MustParse("1474078482417119356")
	err = client.Rest.AddMemberRole(event.GuildID, event.Member.User.ID, memberRoleID)
}

// give member role if user doesnt have it
