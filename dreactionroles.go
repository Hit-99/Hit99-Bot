package main

import (
	"fmt"
	"time"

	devents "github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

// region role reaction message manager

var dRoleChannel snowflake.ID
var dReactionMsgID snowflake.ID

var dRoleReactionIDs = map[string]string{
	"🇳": "1474815372074160370",
	"🇪": "1474815294114369636",
	"🇸": "1474815324280066229",
	"🇦": "1474815257842286612",
	"🇴": "1474815221284474961",
}

func init() {
	dRoleChannel = snowflake.MustParse("1474815023602995290")
	dReactionMsgID = snowflake.MustParse("1499773696477495386")
}

func dInitReactionRoles() {
	// need to stop removing all and then readding. (just check and then remove/add)
	dClient.Rest.RemoveAllReactions(dRoleChannel, dReactionMsgID)
	dAddRoleMsgReactions(dRoleChannel, dReactionMsgID)
}

func dAddRoleMsgReactions(channelID snowflake.ID, messageID snowflake.ID) {
	// add emojis to the role message

	for _, emojiName := range orderedReactionRoles {
		err := dAddReactionToMessage(channelID, messageID, emojiName)
		time.Sleep(300 * time.Millisecond)
		if err != nil {
			fmt.Println("error adding reaction to message:", err)
			return
		}
	}
}

// adds an emoji given channel, message, and emoji ids
func dAddReactionToMessage(channelID snowflake.ID, messageID snowflake.ID, emojiName string) error {
	err := dClient.Rest.AddReaction(channelID, messageID, emojiName)

	return err
}

func dContainsEmoji(emojiSlice map[string]string, emojiName string) bool {
	for emojiSliceName, _ := range emojiSlice {
		if emojiSliceName == emojiName {
			return true
		}
	}
	return false
}

func dUpdateReactionRoles(event *devents.MessageReactionAdd) {

	if event.UserID == dClient.ID() {
		return
	}

	if event.ChannelID != dRoleChannel && event.MessageID != dReactionMsgID {
		return
	}
	if !containsEmoji(dRoleReactionIDs, *event.Emoji.Name) {
		dClient.Rest.RemoveUserReaction(dRoleChannel, dReactionMsgID, *event.Emoji.Name, event.Member.User.ID)
		return

	}

	dToggleRoleForUser(event.Member.User.ID, snowflake.MustParse(dRoleReactionIDs[*event.Emoji.Name]), *event.GuildID)
	dClient.Rest.RemoveUserReaction(dRoleChannel, dReactionMsgID, *event.Emoji.Name, event.Member.User.ID)
	if err != nil {
		fmt.Println("error adding role to user:", err)
		return
	}

}

func dToggleRoleForUser(userID snowflake.ID, roleID snowflake.ID, guildID snowflake.ID) error {

	doesUserHasRole, err := dUserHasRole(guildID, userID, roleID)
	if err != nil {
		return err
	}
	if doesUserHasRole {
		fmt.Println("user has role, removing")
		return dClient.Rest.RemoveMemberRole(guildID, userID, roleID)
	}
	fmt.Println("user does not have role, adding")
	return dClient.Rest.AddMemberRole(guildID, userID, roleID)
}

func dUserHasRole(guidID snowflake.ID, userID snowflake.ID, roleID snowflake.ID) (bool, error) {
	member, err := dClient.Rest.GetMember(guidID, userID)
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
