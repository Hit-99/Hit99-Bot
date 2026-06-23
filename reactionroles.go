package main

import (
	"fmt"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/fluxergo/fluxergo/events"
)

// region role reaction message manager

var roleChannel snowflake.ID
var reactionMsgID snowflake.ID

var orderedReactionRoles = []string{
	"🇳",
	"🇸",
	"🇪",
	"🇦",
	"🇴",
}

var RoleReactionIDs = map[string]string{
	"🇳": "1474171093082538009",
	"🇪": "1474171168777142310",
	"🇸": "1474171228675989527",
	"🇦": "1474171245155393564",
	"🇴": "1474171276428120172",
}

func init() {
	roleChannel = snowflake.MustParse("1474153527752311133")
	reactionMsgID = snowflake.MustParse("1474417595292490005")
}

func initReactionRoles() {
	// need to stop removing all and then readding. (just check and then remove/add)
	fClient.Rest.RemoveAllReactions(roleChannel, reactionMsgID)
	addRoleMsgReactions(roleChannel, reactionMsgID)
}

// doesnt work yet
func addRoleMsgReactions(channelID snowflake.ID, messageID snowflake.ID) {
	// add emojis to the role message

	for _, emojiName := range orderedReactionRoles {
		err := addReactionToMessage(channelID, messageID, emojiName)
		time.Sleep(300 * time.Millisecond)
		if err != nil {
			fmt.Println("error adding reaction to message:", err)
			return
		}
	}
}

func addReactionToMessage(channelID snowflake.ID, messageID snowflake.ID, emojiName string) error {
	err := fClient.Rest.AddReaction(channelID, messageID, emojiName)

	return err
}

func containsEmoji(emojiSlice map[string]string, emojiName string) bool {
	for emojiSliceName, _ := range emojiSlice {
		if emojiSliceName == emojiName {
			return true
		}
	}
	return false
}

func updateReactionRoles(event *events.MessageReactionAdd) {

	// fmt.Println(event.Emoji.String())
	// fmt.Println(event.Emoji.Name)

	if event.UserID == fClient.ID() {
		return
	}

	if event.ChannelID != roleChannel && event.MessageID != reactionMsgID {
		return
	}
	if !containsEmoji(RoleReactionIDs, event.Emoji.Name) {
		fClient.Rest.RemoveUserReaction(roleChannel, reactionMsgID, event.Emoji.Name, event.Member.User.ID)
		return
	}

	toggleRoleForUser(event.Member.User.ID, snowflake.MustParse(RoleReactionIDs[event.Emoji.Name]), *event.GuildID)
	fClient.Rest.RemoveUserReaction(roleChannel, reactionMsgID, event.Emoji.Name, event.Member.User.ID)
	if err != nil {
		fmt.Println("error adding role to user:", err)
		return
	}
}

func toggleRoleForUser(userID snowflake.ID, roleID snowflake.ID, guildID snowflake.ID) error {

	doesUserHasRole, err := userHasRole(guildID, userID, roleID)
	if err != nil {
		return err
	}
	if doesUserHasRole {
		fmt.Println("user has role, removing")
		return fClient.Rest.RemoveMemberRole(guildID, userID, roleID)
	}
	fmt.Println("user does not have role, adding")
	return fClient.Rest.AddMemberRole(guildID, userID, roleID)
}

func userHasRole(guidID snowflake.ID, userID snowflake.ID, roleID snowflake.ID) (bool, error) {
	member, err := fClient.Rest.GetMember(guidID, userID)
	if err != nil {
		return false, err
	}

	for _, r := range member.RoleIDs {
		if r == roleID {
			return true, nil
		}
	}
	return false, nil
}
