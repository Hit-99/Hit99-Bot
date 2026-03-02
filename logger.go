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

func checkAndCreateLogsDir() *os.File {
	err := os.MkdirAll("logs", 0755)
	check(err)
	path := filepath.Join("logs", "latest.log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	check(err)
	return f
}

func formatTime() string {
	currTime := time.Now()
	formattedTime := currTime.Format("15:04:05")
	return formattedTime
}

// logs all fluxer events to log channel and archives them daily

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
		// sends a blank message if its media
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
		logMsgDelete(msgDel)
	}
}

// Reply Log Event
// func userMsgReplyEvent(msgReply *events.)

// Create Role Log Event			(doesnt seem to work)
func roleCreateEvent(roleCreate *events.RoleCreate) {
	roleCreateMsg := fluxer.NewMessageCreate().WithContent(fmt.Sprintf("%s role created", roleCreate.Role))
	logRoleCreate(roleCreate)

	_, err = client.Rest.CreateMessage(logChannelID, roleCreateMsg)
}

// Update Role Log Event			(doesnt seem to work)
func roleUpdateEvent(roleUpdate *events.RoleUpdate) {
	roleUpdateMsg := fluxer.NewMessageCreate().WithContent(fmt.Sprintf("%s role updated", roleUpdate.Role))
	logRoleUpdate(roleUpdate)

	_, err = client.Rest.CreateMessage(logChannelID, roleUpdateMsg)
}

// Delete Role Log Event
func roleDeleteEvent(roleDelete *events.RoleDelete) {
	roleDeleteMsg := fluxer.NewMessageCreate().WithContent(fmt.Sprintf("%s role updated", roleDelete.Role))
	logRoleDelete(roleDelete)

	_, err = client.Rest.CreateMessage(logChannelID, roleDeleteMsg)
}

// Channel Creation Event
func channelCreateEvent(channelCreate *events.GuildChannelCreate) {
	channelCreateMsg := fluxer.NewMessageCreate().WithContent(fmt.Sprintf("%s channel created", channelCreate.ChannelID))
	logChannelCreate(channelCreate)

	_, err = client.Rest.CreateMessage(logChannelID, channelCreateMsg)
}

// Channel Update Event
func channelUpdateEvent(channelUpdate *events.GuildChannelUpdate) {
	channelUpdateMsg := fluxer.NewMessageCreate().WithContent(fmt.Sprintf("%s channel created", channelUpdate.ChannelID))
	logChannelUpdate(channelUpdate)

	_, err = client.Rest.CreateMessage(logChannelID, channelUpdateMsg)
}

// Channel Deletion Event
func channelDeleteEvent(channelDelete *events.GuildChannelDelete) {
	channelDeleteMsg := fluxer.NewMessageCreate().WithContent(fmt.Sprintf("%s channel created", channelDelete.ChannelID))
	logChannelDelete(channelDelete)

	_, err = client.Rest.CreateMessage(logChannelID, channelDeleteMsg)
}

//								-- EVENT LOGGING --

// Write user join event to file
func logUserJoin(join *events.GuildMemberJoin) {
	userJoin := join.Member.User.ID

	f := checkAndCreateLogsDir()
	_, err = fmt.Fprintf(f, "[%s] <%s> joined\n", formatTime(), userJoin)
	check(err)
	defer f.Close()
}

// Write user leave event to file
func logUserLeave(leave *events.GuildMemberLeave) {
	userLeave := leave.Member.User.ID

	f := checkAndCreateLogsDir()
	_, err = fmt.Fprintf(f, "[%s] <%s> left\n", formatTime(), userLeave)
	check(err)
	defer f.Close()
}

// Write message send event to file
func logMsgSend(msgSend *events.GuildMessageCreate) {
	message := msgSend.Message.Content
	messageID := msgSend.MessageID
	// if content = "", dont send (if its media, send attachment name/link)
	author := msgSend.Message.Author.ID
	channelID := msgSend.ChannelID

	f := checkAndCreateLogsDir()
	_, err = fmt.Fprintf(f, "[%s] <%s> <%s> <%s> » %s\n", formatTime(), messageID, channelID, author, message) // [time] <messageID> <channelID> <authorID> » msg
	check(err)
	defer f.Close()

	// logs bot and author if sending message on discord through bridge
}

// Write message edit event to file
func logMsgEdit(msgEdit *events.GuildMessageUpdate) {
	messageNew := msgEdit.Message.Content
	messageOld := msgEdit.OldMessage.Content // returning a blank string
	messageID := msgEdit.MessageID
	// if messageNew = "", dont send
	authorID := msgEdit.Message.Author.ID
	channelID := msgEdit.ChannelID

	f := checkAndCreateLogsDir()
	_, err = fmt.Fprintf(f, "[%s] <%s> <%s> <%s> (BEFORE) » %s   (AFTER) » %s\n", formatTime(), messageID, channelID, authorID, messageOld, messageNew) // [time] <messageID> <channelID> <authorID> (BEFORE) » oldMsg   (AFTER) » newMsg
	check(err)
	defer f.Close()

	// embeds loading and deletion for links counts as an edit update
}

// Write message delete event to file
func logMsgDelete(msgDel *events.GuildMessageDelete) {
	message := msgDel.Message.Content
	messageID := msgDel.MessageID
	authorID := msgDel.Message.Author.ID
	channelID := msgDel.ChannelID

	f := checkAndCreateLogsDir()
	_, err = fmt.Fprintf(f, "[%s] <%s> <%s> <%s> (DELETED) » %s\n", formatTime(), messageID, channelID, authorID, message) // [time] <messageID> <channelID> <authorID> (DELETED) » msg
	check(err)
	defer f.Close()
}

// Write message reply event to file
// need the event listener first

// Write role create event to file
func logRoleCreate(roleCreate *events.RoleCreate) {
	roleID := roleCreate.Role.ID

	f := checkAndCreateLogsDir()
	_, err = fmt.Fprintf(f, "[%s] <%s> role created\n", formatTime(), roleID)
	check(err)
	defer f.Close()
}

// Write role edit event to file
func logRoleUpdate(roleUpdate *events.RoleUpdate) {
	roleID := roleUpdate.Role.ID

	f := checkAndCreateLogsDir()
	_, err = fmt.Fprintf(f, "[%s] <%s> role updated\n", formatTime(), roleID)
	check(err)
	defer f.Close()
}

// Write role delete event to file
func logRoleDelete(roleDelete *events.RoleDelete) {
	roleID := roleDelete.Role.ID

	f := checkAndCreateLogsDir()
	_, err = fmt.Fprintf(f, "[%s] <%s> role deleted\n", formatTime(), roleID)
	check(err)
	defer f.Close()
}

// Write channel create event to file
func logChannelCreate(channelCreate *events.GuildChannelCreate) {
	channelID := channelCreate.ChannelID

	f := checkAndCreateLogsDir()
	_, err = fmt.Fprintf(f, "[%s] <%s> channel created\n", formatTime(), channelID)
	check(err)
	defer f.Close()
}

// Write channel edit event to file
func logChannelUpdate(channelUpdate *events.GuildChannelUpdate) {
	channelID := channelUpdate.ChannelID

	f := checkAndCreateLogsDir()
	_, err = fmt.Fprintf(f, "[%s] <%s> channel updated\n", formatTime(), channelID)
	check(err)
	defer f.Close()
}

// Write channel delete event to file
func logChannelDelete(channelDelete *events.GuildChannelDelete) {
	channelID := channelDelete.ChannelID

	f := checkAndCreateLogsDir()
	_, err = fmt.Fprintf(f, "[%s] <%s> channel deleted\n", formatTime(), channelID)
	check(err)
	defer f.Close()
}
