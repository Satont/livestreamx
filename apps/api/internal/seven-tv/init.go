package seven_tv

import (
	"context"
	"errors"
	"sync"

	"github.com/samber/lo"
	"github.com/satont/stream/apps/api/internal/repositories/user"
)

func (c *SevenTV) init() {
	// wait ws connection
	for {
		if c.wsConn != nil {
			break
		}
	}

	var allUsers []user.User

	ctx, _ := context.WithCancel(context.Background())

	opts := user.FindManyOpts{
		Page:    0,
		PerPage: 1,
	}

	for {
		users, err := c.userRepo.FindMany(ctx, user.FindManyOpts{})
		if err != nil {
			continue
		}

		allUsers = append(allUsers, users.Users...)

		if len(allUsers) < users.Total {
			opts.Page++
		} else {
			break
		}
	}

	for _, user := range allUsers {
		user := user

		go func() {
			if err := c.InitUser(user); err != nil {
				c.logger.Sugar().Error("[7TV] init user", err)
			}
		}()
	}
}

var channelsLock sync.Mutex

func (c *SevenTV) InitUser(user user.User) error {
	for {
		if c.wsConn != nil {
			break
		}
	}

	channelsLock.Lock()
	c.Channels = lo.Filter(
		c.Channels,
		func(item ChannelCache, _ int) bool {
			return item.ChannelID != user.ID
		},
	)
	channelsLock.Unlock()

	if user.SevenTvEmoteSetID == nil {
		return nil
	}

	if err := c.subscribeToEmoteSetUpdates(*user.SevenTvEmoteSetID); err != nil {
		c.logger.Sugar().Error("[7TV] Cannot subscribe to emote set updates", err)
		return err
	}

	c.logger.Sugar().Infow(
		"[7TV] subscribed to emote set updates",
		"user_id", user.ID.String(),
		"emote_set_id", *user.SevenTvEmoteSetID,
	)

	emotes, err := c.fetchEmoteSetEmotes(*user.SevenTvEmoteSetID)
	if err != nil && errors.Is(err, ErrNoEmotes) {
		return err
	}

	c.logger.Sugar().Infow(
		"[7TV] fetched user emotes",
		"user_id", user.ID.String(),
		"emote_set_id", *user.SevenTvEmoteSetID,
		"emotes", len(emotes),
	)

	channelsLock.Lock()
	c.Channels = append(
		c.Channels,
		ChannelCache{
			EmoteSetID: *user.SevenTvEmoteSetID,
			ChannelID:  user.ID,
			Emotes:     emotes,
		},
	)
	channelsLock.Unlock()

	return nil
}
