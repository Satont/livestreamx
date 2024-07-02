package mappers

import (
	"github.com/samber/lo"
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
	message_reaction "github.com/satont/stream/apps/api/internal/repositories/message-reaction"
	seven_tv "github.com/satont/stream/apps/api/internal/seven-tv"
)

func (c *Mapper) DbReactionToGql(
	r message_reaction.MessageReaction,
) gqlmodel.ChatMessageReaction {
	sevenTvChannel, ok := lo.Find(
		c.sevenTv.Channels, func(item seven_tv.ChannelCache) bool {
			return item.ChannelID == r.ChannelID
		},
	)

	sevenTvEmotes := lo.If[[]seven_tv.Emote](!ok, nil).ElseF(
		func() []seven_tv.Emote {
			return lo.Values(sevenTvChannel.Emotes)
		},
	)

	emote, emoteFound := lo.Find(
		sevenTvEmotes,
		func(item seven_tv.Emote) bool {
			return item.Name == r.Reaction
		},
	)

	if emoteFound {
		return gqlmodel.ChatMessageReactionEmote{
			ID:       r.ID.String(),
			Type:     gqlmodel.ChatMessageReactionTypeEmote,
			UserID:   r.UserID,
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
			UserID:    r.UserID,
			Reaction:  r.Reaction,
			MessageID: r.MessageID.String(),
		}
	}
}
