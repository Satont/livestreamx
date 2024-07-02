package seven_tv

type sevenTvEmote struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Flags int    `json:"flags"`
	Data  struct {
		Animated bool `json:"animated"`
		Host     struct {
			URL   string `json:"url"`
			Files []struct {
				Name   string `json:"name"`
				Format string `json:"format"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"files"`
		} `json:"host"`
	} `json:"data"`
}

type SevenTvWebsocketBody struct {
	OldValue *sevenTvEmote `json:"old_value"`
	Value    *sevenTvEmote `json:"value"`
}

const dispatch_opcode = 0

type SevenTvWebsocketPayload struct {
	Op int `json:"op"`
	D  struct {
		Type string `json:"type"`
		Body struct {
			ID      string                 `json:"id"`
			Pulled  []SevenTvWebsocketBody `json:"pulled"`
			Pushed  []SevenTvWebsocketBody `json:"pushed"`
			Updated []SevenTvWebsocketBody `json:"updated"`
			Actor   struct {
				ID          string `json:"id"`
				Username    string `json:"username"`
				DisplayName string `json:"display_name"`
			}
		} `json:"body"`
	}
}
