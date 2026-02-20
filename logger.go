package main

import (
	"fmt"
	"os"

	"github.com/disgoorg/snowflake/v2"
	"github.com/fluxergo/fluxergo/events"
	"github.com/fluxergo/fluxergo/fluxer"
)

var logChannelID snowflake.ID

func init() {
	logChannelID = snowflake.MustParse(os.Getenv("LOG_CHANNEL_ID"))
}

// MAKE THESE EMBEDS

// log all events to a file and/or channel

/*
joins
leaves
messages
replies
*/

// User Join Log Event			(add numbering system)
func userJoinEvent(join *events.GuildMemberJoin) {
	joinLogMsg := fluxer.NewMessageCreate().WithContent(fmt.Sprintf("%s joined Hit99", join.Member))

	_, err = client.Rest.CreateMessage(logChannelID, joinLogMsg)
}

// User Leave Log Event				(doesnt show username when leaving)
func userLeaveEvent(leave *events.GuildMemberLeave) {
	leaveLogMsg := fluxer.NewMessageCreate().WithContent(fmt.Sprintf("%s left Hit99", leave.Member))

	_, err = client.Rest.CreateMessage(logChannelID, leaveLogMsg)
}

// Sent Message Log Event
func userMsgSendEvent(msgSend *events.GuildMessageCreate) {
	if msgSend.Message.Author.ID == client.ID() {
		return
	} else {
		fmt.Printf("%s sent %s", msgSend.Message.Author.Username, msgSend.Message.Content)
	}
}

// Deleted Message Log Event
func userMsgDelEvent(msgDel *events.GuildMessageDelete) {
	if msgDel.Message.Author.ID == client.ID() {
		return
	} else {
		delMsgEmbed(logChannelID, msgDel)
	}
}

// Updated Message Log Event 			(should show a before and after)
func userMsgEditEvent(msgEdit *events.GuildMessageUpdate) {
	if msgEdit.Message.Author.ID == client.ID() {
		return
	} else {
		editMsgEmbed(logChannelID, msgEdit)
	}
}

// Reply Log Event

// Create Role Log Event

// Delete Role Log Event

// Update Role Log Event
