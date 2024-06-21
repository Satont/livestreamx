package seven_tv

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func NewWs(emoteSetID string) error {
	u := url.URL{Scheme: "wss", Host: "events.7tv.io", Path: "/v3"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}

	if err := conn.WriteMessage(
		websocket.TextMessage,
		[]byte(
			fmt.Sprintf(
				`{"op":35,"d":{"type":"emote_set.update","condition":{"object_id":"%s"}}}`,
				emoteSetID,
			),
		),
	); err != nil {
		panic(err)
	}

	go func() {
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	return nil
}
