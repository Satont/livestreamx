package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
	chat_message "github.com/satont/stream/apps/api/internal/repositories/chat-message"
	chat_messages_with_user "github.com/satont/stream/apps/api/internal/repositories/chat-messages-with-user"
)

// SendMessage is the resolver for the sendMessage field.
func (r *mutationResolver) SendMessage(ctx context.Context, input gqlmodel.SendMessageInput) (bool, error) {
	userId, err := r.sessionStorage.GetUserID(ctx)
	if err != nil {
		return false, err
	}

	text := input.Text
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, "\r", "")

	message, err := r.chatMessageRepo.Create(
		ctx,
		chat_message.CreateChatMessageOpts{
			SenderID: uuid.MustParse(userId),
			Text:     text,
		},
	)
	if err != nil {
		return false, fmt.Errorf("failed to create chat message: %w", err)
	}

	messageWithUser, err := r.chatMessagesWithUserRepo.FindByID(ctx, message.ID)
	if err != nil {
		return false, fmt.Errorf("failed to find chat message: %w", err)
	}

	gqlMessage := r.converter.ChatMessageWithUser(ctx, messageWithUser)

	for _, ch := range r.chatListenersChannels {
		ch <- &gqlMessage
	}

	return true, nil
}

// ChatMessagesLatest is the resolver for the chatMessagesLatest field.
func (r *queryResolver) ChatMessagesLatest(ctx context.Context, limit *int) ([]gqlmodel.ChatMessage, error) {
	limitQuery := 100
	if limit != nil {
		limitQuery = *limit
	}

	messages, err := r.chatMessagesWithUserRepo.FindLatest(
		ctx, chat_messages_with_user.FindLatestOpts{
			Limit: limitQuery,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find chat messages: %w", err)
	}

	gqlMessages := make([]gqlmodel.ChatMessage, 0, len(messages))
	for _, m := range messages {
		gqlMessages = append(
			gqlMessages,
			r.converter.ChatMessageWithUser(ctx, &m),
		)
	}

	return gqlMessages, nil
}

// ChatMessages is the resolver for the chatMessages field.
func (r *subscriptionResolver) ChatMessages(ctx context.Context) (<-chan *gqlmodel.ChatMessage, error) {
	id := uuid.NewString()

	r.chatListenersChannels[id] = make(chan *gqlmodel.ChatMessage, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(r.chatListenersChannels[id])
				delete(r.chatListenersChannels, id)
				return
			}
		}
	}()

	return r.chatListenersChannels[id], nil
}
