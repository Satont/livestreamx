package mappers

import (
	"context"
	"regexp"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
	chat_message "github.com/satont/stream/apps/api/internal/repositories/chat-message"
	seven_tv "github.com/satont/stream/apps/api/internal/seven-tv"
)

var linkRegexp = regexp.MustCompile(`(https?://[^\s]+)`)
var mentionRegexp = regexp.MustCompile(`@([a-zA-Z0-9_]+)`)

func (c *Mapper) ChatMessageWithUser(
	ctx context.Context,
	m *chat_message.Message,
) gqlmodel.ChatMessage {
	splittedText := strings.Fields(m.Text)
	var segments []gqlmodel.MessageSegment

	emotesSlice := lo.Values(c.sevenTv.Emotes)

	for _, text := range splittedText {
		emote, emoteFound := lo.Find(
			emotesSlice, func(item seven_tv.Emote) bool {
				return item.Name == text
			},
		)

		if emoteFound {
			segments = append(
				segments,
				gqlmodel.MessageSegmentEmote{
					Content: text,
					Type:    gqlmodel.MessageSegmentTypeEmote,
					Emote: &gqlmodel.Emote{
						ID:     emote.ID,
						Name:   emote.Name,
						URL:    emote.URL,
						Width:  emote.Width,
						Height: emote.Height,
					},
				},
			)
		} else if linkRegexp.MatchString(text) {
			segments = append(
				segments, gqlmodel.MessageSegmentLink{
					Content: text,
					Type:    gqlmodel.MessageSegmentTypeLink,
				},
			)
		} else if mentionRegexp.MatchString(text) {
			mentionSegment := gqlmodel.MessageSegmentMention{
				Content: text,
				Type:    gqlmodel.MessageSegmentTypeMention,
			}

			user, err := c.userRepo.FindByName(ctx, text[1:])
			if err == nil && user != nil {
				mentionSegment.UserID = user.ID
				segments = append(segments, mentionSegment)
			} else {
				segments = append(
					segments,
					gqlmodel.MessageSegmentText{
						Content: text,
						Type:    gqlmodel.MessageSegmentTypeText,
					},
				)
			}
		} else {
			segments = append(
				segments, gqlmodel.MessageSegmentText{
					Content: text,
					Type:    gqlmodel.MessageSegmentTypeText,
				},
			)
		}
	}

	reactions := make([]gqlmodel.ChatMessageReaction, 0, len(m.Reactions))
	for _, r := range m.Reactions {
		reactions = append(reactions, c.DbReactionToGql(r))
	}

	return gqlmodel.ChatMessage{
		ID:        m.ID.String(),
		ChannelID: m.ChannelID,
		SenderID:  m.SenderID,
		Segments:  segments,
		CreatedAt: m.CreatedAt,
		Reactions: reactions,
		ReplyTo:   m.ReplyTo,
	}
}
