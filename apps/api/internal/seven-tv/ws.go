package seven_tv

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
)

func (c *SevenTV) openWebSocket() {
	u := url.URL{Scheme: "wss", Host: "events.7tv.io", Path: "/v3"}

	go func() {
		for {
			conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			defer conn.Close()
			if err != nil {
				fmt.Println(err)
				continue
			}

			if err := conn.WriteMessage(
				websocket.TextMessage,
				[]byte(
					fmt.Sprintf(
						`{"op":35,"d":{"type":"emote_set.update","condition":{"object_id":"%s"}}}`,
						c.config.SevenTVEmoteSetID,
					),
				),
			); err != nil {
				fmt.Println(err)
				continue
			}

		readLoop:
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					break readLoop
				}

				data := SevenTvWebsocketPayload{}
				if err := json.Unmarshal(message, &data); err != nil {
					log.Println("unmarshal:", err)
					continue
				}

				if data.Op != dispatch_opcode {
					continue
				}

				for _, emote := range data.D.Body.Pushed {
					log.Printf("added: %s", emote.Value.Name)

					c.Emotes[emote.Value.ID] = Emote{
						ID:     emote.Value.ID,
						Name:   emote.Value.Name,
						URL:    fmt.Sprintf("%s/%s", emote.Value.Data.Host.URL, "1x.webp"),
						Width:  emote.Value.Data.Host.Files[0].Width,
						Height: emote.Value.Data.Host.Files[0].Height,
					}
				}

				for _, pulledBody := range data.D.Body.Pulled {
					if pulledBody.OldValue != nil && pulledBody.Value == nil {
						log.Printf("deleted: %s", pulledBody.OldValue.Name)
						delete(c.Emotes, pulledBody.OldValue.ID)
					}
				}

				for _, emote := range data.D.Body.Updated {
					log.Printf("updated: %s", emote.Value.Name)
				}
			}

			time.Sleep(500 * time.Millisecond)
		}
	}()
}
