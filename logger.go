package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/fluxergo/fluxergo/events"
	"github.com/fluxergo/fluxergo/fluxer"
)

var logChannelID snowflake.ID

func init() {
	logChannelID = snowflake.MustParse(os.Getenv("LOG_CHANNEL_ID"))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkForLogsDir() {
	err := os.MkdirAll("logs", 0755)
	check(err)
}

// log all messages and replies to file
// log other events to log channel (joins, leaves, edits)

// User Join Log Event			(add numbering system)
func userJoinEvent(join *events.GuildMemberJoin) {
	joinLogMsg := fluxer.NewMessageCreate().WithContent(fmt.Sprintf("%s joined Hit99", join.Member))
	logUserJoin(join)

	_, err = client.Rest.CreateMessage(logChannelID, joinLogMsg)
}

// User Leave Log Event				(doesnt show username when leaving)
func userLeaveEvent(leave *events.GuildMemberLeave) {
	leaveLogMsg := fluxer.NewMessageCreate().WithContent(fmt.Sprintf("%s left Hit99", leave.Member))
	logUserLeave(leave)

	_, err = client.Rest.CreateMessage(logChannelID, leaveLogMsg)
}

// Sent Message Log Event
func userMsgSendEvent(msgSend *events.GuildMessageCreate) {
	if msgSend.Message.Author.ID == client.ID() {
		return
	} else {
		logMsgSend(msgSend)
	}
}

// Edit Message Log Event 			(should show a before and after)
func userMsgEditEvent(msgEdit *events.GuildMessageUpdate) {
	if msgEdit.Message.Author.ID == client.ID() {
		return
	} else {
		editMsgEmbed(logChannelID, msgEdit)
		logMsgEdit(msgEdit)
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

// Reply Log Event
// func userMsgReplyEvent(msgReply *events.)

// Create Role Log Event
func roleCreateEvent(roleCreate *events.RoleCreate) {
	roleCreateMsg := fluxer.NewMessageCreate().WithContent(fmt.Sprintf("%s role created", roleCreate.Role))
	logRoleCreate(roleCreate)

	_, err = client.Rest.CreateMessage(logChannelID, roleCreateMsg)
}

// Update Role Log Event

// Delete Role Log Event

// Channel Creation Event

// Channel Edit Event

// Channel Deletion Event

//

// create .log file
// write all log events to file (msg send, edit, delete, replies, and media) (role creation, edits, and deletes) (channel creation, edits, and deletes)
// save media into folder
// add timestamps
// archive daily into .gz file

// Write user join event to file
func logUserJoin(join *events.GuildMemberJoin) {
	currTime := time.Now()
	formattedTime := currTime.Format("15:04:05")

	userJoin := join.Member.User.ID

	checkForLogsDir()
	path := filepath.Join("logs", "latest.log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	check(err)
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("[%s] <%s> joined\n", formattedTime, userJoin))
	check(err)
}

// Write user leave event to file
func logUserLeave(leave *events.GuildMemberLeave) {
	currTime := time.Now()
	formattedTime := currTime.Format("15:04:05")

	userLeave := leave.Member.User.ID

	checkForLogsDir()
	path := filepath.Join("logs", "latest.log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	check(err)
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("[%s] <%s> left\n", formattedTime, userLeave))
	check(err)
}

// Write message send event to file
func logMsgSend(msgSend *events.GuildMessageCreate) {
	currTime := time.Now()
	formattedTime := currTime.Format("15:04:05")

	message := msgSend.Message.Content
	// if content = "", dont send
	author := msgSend.Message.Author.ID
	channelID := msgSend.ChannelID

	checkForLogsDir()
	path := filepath.Join("logs", "latest.log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	check(err)
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("[%s] <%s> <%s> » %s\n", formattedTime, channelID, author, message))
	check(err)

	// logs bot and author if sending message on discord through bridge
}

// Write message edit event to file
func logMsgEdit(msgEdit *events.GuildMessageUpdate) {
	currTime := time.Now()
	formattedTime := currTime.Format("15:04:05")

	messageNew := msgEdit.Message.Content
	messageOld := msgEdit.OldMessage.Content
	// if content = "", dont send
	author := msgEdit.Message.Author.ID
	channelID := msgEdit.ChannelID

	checkForLogsDir()
	path := filepath.Join("logs", "latest.log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	check(err)
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("[%s] <%s> <%s> (BEFORE) » %s   (AFTER) » %s\n", formattedTime, channelID, author, messageOld, messageNew))
	check(err)

	// embeds loading for links counts as an edit update
}

// Write message leave event to file

// Write message reply event to file

// Write role create event to file
func logRoleCreate(roleCreate *events.RoleCreate) {
	currTime := time.Now()
	formattedTime := currTime.Format("15:04:05")

	roleID := roleCreate.Role.ID

	checkForLogsDir()
	path := filepath.Join("logs", "latest.log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	check(err)
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("[%s] <%s> role create\n", formattedTime, roleID))
	check(err)
}

// Write role edit event to file

// Write role delete event to file

// Write channel create event to file

// Write channel edit event to file

// Write channel delete event to file
