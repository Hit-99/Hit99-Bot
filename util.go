package main

import (
	"github.com/disgoorg/snowflake/v2"
	"github.com/fluxergo/fluxergo/fluxer"
)

func getMessageById(channelID snowflake.ID, messageID snowflake.ID) (*fluxer.Message, error) {
	msg, err := client.Rest.GetMessage(channelID, messageID)
	return msg, err
}
