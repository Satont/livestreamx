package seven_tv

import (
	"errors"
	"fmt"
	"time"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
)

type fetchEmoteSetResponseFile struct {
	Name       string `json:"name"`
	StaticName string `json:"static_name"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	FrameCount int    `json:"frame_count"`
	Size       int    `json:"size"`
	Format     string `json:"format"`
}

type fetchEmoteSetResponse struct {
	Id         string        `json:"id"`
	Name       string        `json:"name"`
	Flags      int           `json:"flags"`
	Tags       []interface{} `json:"tags"`
	Immutable  bool          `json:"immutable"`
	Privileged bool          `json:"privileged"`
	Emotes     []struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		Flags     int    `json:"flags"`
		Timestamp int64  `json:"timestamp"`
		ActorId   string `json:"actor_id"`
		Data      struct {
			Id        string   `json:"id"`
			Name      string   `json:"name"`
			Flags     int      `json:"flags"`
			Lifecycle int      `json:"lifecycle"`
			State     []string `json:"state"`
			Listed    bool     `json:"listed"`
			Animated  bool     `json:"animated"`
			Owner     struct {
				Id          string `json:"id"`
				Username    string `json:"username"`
				DisplayName string `json:"display_name"`
				AvatarUrl   string `json:"avatar_url"`
				Style       struct {
				} `json:"style"`
				Roles []string `json:"roles"`
			} `json:"owner"`
			Host struct {
				Url   string                      `json:"url"`
				Files []fetchEmoteSetResponseFile `json:"files"`
			} `json:"host"`
		} `json:"data"`
	} `json:"emotes"`
	EmoteCount int `json:"emote_count"`
	Capacity   int `json:"capacity"`
	Owner      struct {
		Id          string `json:"id"`
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		AvatarUrl   string `json:"avatar_url"`
		Style       struct {
		} `json:"style"`
		Roles []string `json:"roles"`
	} `json:"owner"`
}

var ErrNoEmotes = errors.New("no emotes")

func (c *SevenTV) fetchEmoteSetEmotes(emoteSetID string) (map[string]Emote, error) {
	data := &fetchEmoteSetResponse{}

	_, err := req.
		SetSuccessResult(&data).
		SetRetryCount(20).
		SetRetryFixedInterval(200 * time.Millisecond).
		Get(
			fmt.Sprintf(
				"https://7tv.io/v3/emote-sets/%s?t=%v",
				emoteSetID,
				time.Now().UnixMilli(),
			),
		)
	if err != nil {
		return nil, err
	}

	if len(data.Emotes) == 0 {
		return nil, ErrNoEmotes
	}

	emotes := make(map[string]Emote)

	for _, emote := range data.Emotes {
		file, _ := lo.Find(
			emote.Data.Host.Files,
			func(item fetchEmoteSetResponseFile) bool {
				return item.Name == "1x.webp"
			},
		)

		emotes[emote.Id] = Emote{
			ID:     emote.Id,
			Name:   emote.Name,
			URL:    fmt.Sprintf("%s/%s", emote.Data.Host.Url, "1x.webp"),
			Width:  file.Width,
			Height: file.Height,
		}
	}

	return emotes, nil
}
