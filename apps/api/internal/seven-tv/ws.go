package seven_tv

import (
	"fmt"
	"net/url"
	"time"

	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
	system_messages "github.com/satont/stream/apps/api/internal/system-messages"
)

func (c *SevenTV) openWebSocket() {
	u := url.URL{Scheme: "wss", Host: "events.7tv.io", Path: "/v3"}

	go func() {
		for {
			conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			defer conn.Close()
			if err != nil {
				c.logger.Sugar().Error("[7TV] dial", err)
				continue
			}

			c.logger.Sugar().Info("[7TV] socket connected")
			c.wsConn = conn

		readLoop:
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					c.logger.Sugar().Error("[7TV] read", err)
					break readLoop
				}

				data := SevenTvWebsocketPayload{}
				if err := json.Unmarshal(message, &data); err != nil {
					c.logger.Sugar().Error("[7TV] unmarshal", err, string(message))
					continue
				}

				if data.Op != dispatch_opcode {
					continue
				}

				for _, emote := range data.D.Body.Pushed {
					for _, ch := range c.Channels {
						if ch.EmoteSetID != data.D.Body.ID {
							continue
						}

						c.logger.Sugar().Infow(
							"[7TV] added",
							"emote_name",
							emote.Value.Name,
							"emote_set_id",
							data.D.Body.ID,
							"channel_id",
							ch.ChannelID,
						)

						ch.Emotes[emote.Value.ID] = Emote{
							ID:     emote.Value.ID,
							Name:   emote.Value.Name,
							URL:    fmt.Sprintf("%s/%s", emote.Value.Data.Host.URL, "1x.webp"),
							Width:  emote.Value.Data.Host.Files[0].Width,
							Height: emote.Value.Data.Host.Files[0].Height,
						}

						c.subscriptionsRouter.Publish(
							system_messages.BuildSubscriptionEmoteAddedKey(ch.ChannelID),
							&gqlmodel.SystemMessageEmoteAdded{
								Type:      gqlmodel.SystemMessageTypeEmoteAdded,
								CreatedAt: time.Now().UTC(),
								Emote: &gqlmodel.Emote{
									ID:     emote.Value.ID,
									Name:   emote.Value.Name,
									URL:    fmt.Sprintf("%s/%s", emote.Value.Data.Host.URL, "1x.webp"),
									Width:  emote.Value.Data.Host.Files[0].Width,
									Height: emote.Value.Data.Host.Files[0].Height,
								},
								Actor: &gqlmodel.SystemMessageEmoteActor{
									ID:          data.D.Body.Actor.ID,
									Name:        data.D.Body.Actor.Username,
									DisplayName: data.D.Body.Actor.DisplayName,
								},
							},
						)
					}
				}

				for _, pulledBody := range data.D.Body.Pulled {
					if pulledBody.OldValue != nil && pulledBody.Value == nil {
						for _, ch := range c.Channels {
							if ch.EmoteSetID != data.D.Body.ID {
								continue
							}

							c.logger.Sugar().Infow(
								"[7TV] deleted",
								"emote_name",
								pulledBody.OldValue.Name,
								"emote_set_id",
								data.D.Body.ID,
								"channel_id",
								ch.ChannelID,
							)

							emoteData, err := c.fetchSingleEmote(pulledBody.OldValue.ID)
							if err != nil {
								c.logger.Sugar().Error("[7TV] fetch single emote", err)
								continue
							}

							c.subscriptionsRouter.Publish(
								system_messages.BuildSubscriptionEmoteRemovedKey(ch.ChannelID),
								&gqlmodel.SystemMessageEmoteRemoved{
									Type:      gqlmodel.SystemMessageTypeEmoteRemoved,
									CreatedAt: time.Now().UTC(),
									Emote: &gqlmodel.Emote{
										ID:     emoteData.Id,
										Name:   emoteData.Name,
										URL:    fmt.Sprintf("%s/%s", emoteData.Host.Url, "1x.webp"),
										Width:  emoteData.Host.Files[0].Width,
										Height: emoteData.Host.Files[0].Height,
									},
									Actor: &gqlmodel.SystemMessageEmoteActor{
										ID:          data.D.Body.Actor.ID,
										Name:        data.D.Body.Actor.Username,
										DisplayName: data.D.Body.Actor.DisplayName,
									},
								},
							)

							delete(ch.Emotes, pulledBody.OldValue.ID)
						}
					}
				}

				for _, emote := range data.D.Body.Updated {
					for _, ch := range c.Channels {
						if ch.EmoteSetID != data.D.Body.ID {
							continue
						}

						if _, ok := ch.Emotes[emote.Value.ID]; ok {
							ch.Emotes[emote.Value.ID] = Emote{
								ID:     emote.Value.ID,
								Name:   emote.Value.Name,
								URL:    fmt.Sprintf("%s/%s", emote.Value.Data.Host.URL, "1x.webp"),
								Width:  emote.Value.Data.Host.Files[0].Width,
								Height: emote.Value.Data.Host.Files[0].Height,
							}
						}

						c.logger.Sugar().Infow(
							"[7TV] updated",
							"emote_name",
							emote.Value.Name,
							"emote_set_id",
							data.D.Body.ID,
						)
					}
				}
			}

			time.Sleep(500 * time.Millisecond)
			continue
		}
	}()
}

func (c *SevenTV) subscribeToEmoteSetUpdates(emoteSetID string) error {
	// do not subscribe to the same emote set
	for _, ch := range c.Channels {
		if ch.EmoteSetID == emoteSetID {
			return nil
		}
	}

	return c.wsConn.WriteMessage(
		websocket.TextMessage,
		[]byte(
			fmt.Sprintf(
				`{"op":35,"d":{"type":"emote_set.update","condition":{"object_id":"%v"}}}`,
				emoteSetID,
			),
		),
	)
}
