package system_messages

import (
	"github.com/google/uuid"
)

func BuildSubscriptionEmoteAddedKey(channelID uuid.UUID) string {
	return "emote_added:" + channelID.String()
}

func BuildSubscriptionEmoteRemovedKey(channelID uuid.UUID) string {
	return "emote_removed:" + channelID.String()
}

func BuildUserJoinedChannelKey(channelID uuid.UUID) string {
	return "user_joined_channel:" + channelID.String()
}
