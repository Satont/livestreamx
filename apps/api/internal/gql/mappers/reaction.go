package mappers

import (
	"github.com/samber/lo"
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
	message_reaction "github.com/satont/stream/apps/api/internal/repositories/message-reaction"
	seven_tv "github.com/satont/stream/apps/api/internal/seven-tv"
)

func (c *Converters) DbReactionToGql(
	r message_reaction.MessageReaction,
	user *gqlmodel.User,
) gqlmodel.ChatMessageReaction {
	emotesSlice := lo.Values(c.sevenTv.Emotes)

	emote, emoteFound := lo.Find(
		emotesSlice, func(item seven_tv.Emote) bool {
			return item.Name == r.Reaction
		},
	)

	if emoteFound {
		return gqlmodel.ChatMessageReactionEmote{
			ID:       r.ID.String(),
			Type:     gqlmodel.ChatMessageReactionTypeEmote,
			User:     user,
			Reaction: r.Reaction,
			Emote: &gqlmodel.Emote{
				ID:     emote.ID,
				Name:   emote.Name,
				URL:    emote.URL,
				Width:  emote.Width,
				Height: emote.Height,
			},
			MessageID: r.MessageID.String(),
		}
	} else {
		return gqlmodel.ChatMessageReactionEmoji{
			ID:        r.ID.String(),
			Type:      gqlmodel.ChatMessageReactionTypeEmoji,
			User:      user,
			Reaction:  r.Reaction,
			MessageID: r.MessageID.String(),
		}
	}
}