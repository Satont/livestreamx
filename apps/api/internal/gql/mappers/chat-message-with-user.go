package mappers

import (
	"context"
	"regexp"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
	chat_messages_with_user "github.com/satont/stream/apps/api/internal/repositories/chat-messages-with-user"
	seven_tv "github.com/satont/stream/apps/api/internal/seven-tv"
)

var linkRegexp = regexp.MustCompile(`(https?://[^\s]+)`)
var mentionRegexp = regexp.MustCompile(`@([a-zA-Z0-9_]+)`)

func (c *Converters) ChatMessageWithUser(
	ctx context.Context,
	m *chat_messages_with_user.MessageWithUser,
) gqlmodel.ChatMessage {
	splittedText := strings.Fields(m.Message.Text)
	var segments []gqlmodel.MessageSegment
	usersDbCache := make(map[string]*gqlmodel.User)

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

			if cachedUser, ok := usersDbCache[text[1:]]; ok {
				mentionSegment.User = cachedUser
				segments = append(segments, mentionSegment)
				continue
			}

			user, err := c.userRepo.FindByName(ctx, text[1:])
			if err == nil && user != nil {
				userGql := c.DbUserToGql(*user)

				usersDbCache[text[1:]] = &userGql
				mentionSegment.User = &userGql

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
		var user *gqlmodel.User
		if u, ok := usersDbCache[r.UserID.String()]; ok {
			user = u
		} else {
			u, err := c.userRepo.FindByID(ctx, r.UserID)
			if err == nil && u != nil {
				userGql := c.DbUserToGql(*u)

				usersDbCache[r.UserID.String()] = &userGql
				user = &userGql
			}
		}

		reactions = append(reactions, c.DbReactionToGql(r, user))
	}

	sender := c.DbUserToGql(m.User)

	return gqlmodel.ChatMessage{
		ID:        m.Message.ID.String(),
		Segments:  segments,
		Sender:    &sender,
		CreatedAt: m.Message.CreatedAt,
		Reactions: reactions,
	}
}
