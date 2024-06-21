package converters

import (
	"context"
	"regexp"
	"strings"

	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
	chat_messages_with_user "github.com/satont/stream/apps/api/internal/repositories/chat-messages-with-user"
)

var linkRegexp = regexp.MustCompile(`(https?://[^\s]+)`)
var mentionRegexp = regexp.MustCompile(`@([a-zA-Z0-9_]+)`)

func (c *Converters) ChatMessageWithUser(
	ctx context.Context,
	m *chat_messages_with_user.MessageWithUser,
) gqlmodel.ChatMessage {
	splittedText := strings.Fields(m.Message.Text)
	var segments []gqlmodel.MessageSegment
	for _, text := range splittedText {
		if linkRegexp.MatchString(text) {
			segments = append(
				segments, gqlmodel.MessageSegmentLink{
					Content: text,
					Type:    gqlmodel.MessageSegmentTypeLink,
				},
			)
		} else if mentionRegexp.MatchString(text) {
			user, err := c.userRepo.FindByName(ctx, text[1:])
			if err == nil && user != nil {
				segments = append(
					segments,
					gqlmodel.MessageSegmentMention{
						Content: text,
						Type:    gqlmodel.MessageSegmentTypeMention,
						User: &gqlmodel.User{
							ID:          user.ID.String(),
							Name:        user.Name,
							DisplayName: user.DisplayName,
							Color:       user.Color,
							Roles:       nil,
							IsBanned:    false,
							CreatedAt:   user.CreatedAt,
							AvatarURL:   user.AvatarUrl,
						},
					},
				)
			} else {
				segments = append(
					segments, gqlmodel.MessageSegmentText{
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

	return gqlmodel.ChatMessage{
		ID:       m.Message.ID.String(),
		Segments: segments,
		Sender: &gqlmodel.User{
			ID:          m.User.ID.String(),
			Name:        m.User.Name,
			DisplayName: m.User.DisplayName,
			Color:       m.User.Color,
			Roles:       nil,
			IsBanned:    false,
			CreatedAt:   m.User.CreatedAt,
			AvatarURL:   m.User.AvatarUrl,
		},
		CreatedAt: m.Message.CreatedAt,
	}
}
