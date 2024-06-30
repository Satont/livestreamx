package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/99designs/gqlgen/graphql"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
	"github.com/satont/stream/apps/api/internal/httpserver/middlewares"
	chat_message "github.com/satont/stream/apps/api/internal/repositories/chat-message"
	chat_messages_with_user "github.com/satont/stream/apps/api/internal/repositories/chat-messages-with-user"
	message_reaction "github.com/satont/stream/apps/api/internal/repositories/message-reaction"
)

// SendMessage is the resolver for the sendMessage field.
func (r *mutationResolver) SendMessage(ctx context.Context, input gqlmodel.SendMessageInput) (
	bool,
	error,
) {
	userId, err := r.sessionStorage.GetUserID(ctx)
	if err != nil {
		return false, err
	}

	text := input.Text
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, "\r", "")

	if utf8.RuneCountInString(text) > 700 {
		return false, nil
	}

	if badSymbolsRegexp.MatchString(text) {
		return false, nil
	}

	message, err := r.chatMessageRepo.Create(
		ctx,
		chat_message.CreateChatMessageOpts{
			SenderID:  uuid.MustParse(userId),
			ChannelID: input.ChannelID,
			Text:      text,
			ReplyTo:   input.ReplyTo.Value(),
		},
	)
	if err != nil {
		return false, fmt.Errorf("failed to create chat message: %w", err)
	}

	messageWithUser, err := r.chatMessagesWithUserRepo.FindByID(ctx, message.ID)
	if err != nil {
		return false, fmt.Errorf("failed to find chat message: %w", err)
	}

	go func() {
		gqlMessage := r.mapper.ChatMessageWithUser(ctx, messageWithUser)
		if err := r.subscriptionRouter.Publish(
			fmt.Sprintf(
				chatMessagesSubscriptionKey,
				gqlMessage.ChannelID,
			),
			&gqlMessage,
		); err != nil {
			fmt.Println(err)
		}
	}()

	return true, nil
}

// AttachFile is the resolver for the attachFile field.
func (r *mutationResolver) AttachFile(
	ctx context.Context,
	file graphql.Upload,
) (*gqlmodel.AttachedFile, error) {
	_, err := r.s3.PutObject(
		ctx,
		r.config.S3Bucket,
		fmt.Sprintf("badges/%s", file.Filename),
		file.File,
		file.Size,
		minio.PutObjectOptions{
			ContentType: file.ContentType,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to put object to s3: %w", err)
	}

	// attachedFile := &gqlmodel.AttachedFile{
	// 	ID:        "",
	// 	URL:       "",
	// 	Name:      "",
	// 	Size:      int(uploadInfo.Size),
	// 	MimeType:  file.ContentType,
	// 	CreatedAt: time.Time{},
	// }

	panic(fmt.Errorf("not implemented: AttachFile - attachFile"))
}

// AddReaction is the resolver for the addReaction field.
func (r *mutationResolver) AddReaction(
	ctx context.Context,
	messageID string,
	content string,
	channelID uuid.UUID,
) (bool, error) {
	userID, err := r.sessionStorage.GetUserID(ctx)
	if err != nil {
		return false, err
	}

	user := middlewares.GetUserFromContext(ctx)
	if user == nil {
		return false, fmt.Errorf("user not found")
	}

	newReaction, err := r.messageReactionRepo.Create(
		ctx,
		message_reaction.CreateMessageReactionOpts{
			MessageID: uuid.MustParse(messageID),
			UserID:    uuid.MustParse(userID),
			Reaction:  content,
		},
	)
	if err != nil {
		return false, err
	}

	go func() {
		gqlUser := r.mapper.DbUserToGql(*user)
		gqlReaction := r.mapper.DbReactionToGql(*newReaction, &gqlUser)

		// TODO: reaction should use channel id for publish and subscription
		if err := r.subscriptionRouter.Publish(
			fmt.Sprintf(
				chatMessageReactionKey,
				channelID.String(),
			),
			&gqlReaction,
		); err != nil {
			fmt.Println(err)
		}
	}()

	return true, nil
}

// ChatMessagesLatest is the resolver for the chatMessagesLatest field.
func (r *queryResolver) ChatMessagesLatest(
	ctx context.Context,
	channelID uuid.UUID,
	limit *int,
) ([]gqlmodel.ChatMessage, error) {
	limitQuery := 100
	if limit != nil {
		limitQuery = *limit
	}

	messages, err := r.chatMessagesWithUserRepo.FindLatest(
		ctx, chat_messages_with_user.FindLatestOpts{
			Limit:     limitQuery,
			ChannelID: channelID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find chat messages: %w", err)
	}

	gqlMessages := make([]gqlmodel.ChatMessage, 0, len(messages))
	for _, m := range messages {
		gqlMessages = append(
			gqlMessages,
			r.mapper.ChatMessageWithUser(ctx, &m),
		)
	}

	slices.SortFunc(
		gqlMessages,
		func(a, b gqlmodel.ChatMessage) int {
			if a.CreatedAt.UnixMilli() > b.CreatedAt.UnixMilli() {
				return 1
			} else {
				return -1
			}
		},
	)

	return gqlMessages, nil
}

// GetEmotes is the resolver for the getEmotes field.
func (r *queryResolver) GetEmotes(ctx context.Context, channelID uuid.UUID) (
	[]gqlmodel.Emote,
	error,
) {
	emotes := make([]gqlmodel.Emote, 0, len(r.sevenTv.Emotes))
	for _, emote := range r.sevenTv.Emotes {
		emotes = append(
			emotes,
			gqlmodel.Emote{
				ID:     emote.ID,
				Name:   emote.Name,
				URL:    emote.URL,
				Width:  emote.Width,
				Height: emote.Height,
			},
		)
	}

	slices.SortFunc(
		emotes,
		func(a, b gqlmodel.Emote) int {
			return strings.Compare(
				strings.ToLower(a.Name),
				strings.ToLower(b.Name),
			)
		},
	)

	return emotes, nil
}

// ChatMessages is the resolver for the chatMessages field.
func (r *subscriptionResolver) ChatMessages(
	ctx context.Context,
	channelID uuid.UUID,
) (<-chan *gqlmodel.ChatMessage, error) {
	channel := make(chan *gqlmodel.ChatMessage, 1)

	go func() {
		sub, err := r.subscriptionRouter.Subscribe(
			[]string{
				fmt.Sprintf(
					chatMessagesSubscriptionKey,
					channelID,
				),
			},
		)
		if err != nil {
			panic(err)
		}
		defer func() {
			sub.Unsubscribe()
			close(channel)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case data := <-sub.GetChannel():
				var message gqlmodel.ChatMessage
				if err := json.Unmarshal(data, &message); err != nil {
					panic(err)
				}

				channel <- &message
			}
		}
	}()

	return channel, nil
}

// SystemMessages is the resolver for the systemMessages field.
func (r *subscriptionResolver) SystemMessages(
	ctx context.Context,
	channelID uuid.UUID,
) (<-chan gqlmodel.SystemMessage, error) {
	panic(fmt.Errorf("not implemented: SystemMessages - systemMessages"))
}

// ReactionAdd is the resolver for the reactionAdd field.
func (r *subscriptionResolver) ReactionAdd(
	ctx context.Context,
	channelID uuid.UUID,
) (<-chan gqlmodel.ChatMessageReaction, error) {
	channel := make(chan gqlmodel.ChatMessageReaction, 1)

	go func() {
		sub, err := r.subscriptionRouter.Subscribe(
			[]string{
				fmt.Sprintf(
					chatMessageReactionKey,
					channelID,
				),
			},
		)
		if err != nil {
			panic(err)
		}
		defer func() {
			sub.Unsubscribe()
			close(channel)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case data := <-sub.GetChannel():
				var notificationRaw map[string]any
				if err := json.Unmarshal(data, &notificationRaw); err != nil {
					r.logger.Sugar().Error(err)
					continue
				}

				t, ok := notificationRaw["type"].(string)
				if !ok {
					r.logger.Sugar().Error("failed to get notification type")
					continue
				}

				switch gqlmodel.ChatMessageReactionType(t) {
				case gqlmodel.ChatMessageReactionTypeEmote:
					var notification gqlmodel.ChatMessageReactionEmote
					if err := json.Unmarshal(data, &notification); err != nil {
						r.logger.Sugar().Error(err)
						continue
					}

					channel <- notification
				case gqlmodel.ChatMessageReactionTypeEmoji:
					var notification gqlmodel.ChatMessageReactionEmoji
					if err := json.Unmarshal(data, &notification); err != nil {
						r.logger.Sugar().Error(err)
						continue
					}

					channel <- notification
				}
			}
		}
	}()

	return channel, nil
}
