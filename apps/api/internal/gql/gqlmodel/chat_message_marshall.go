package gqlmodel

import (
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

func (cm *ChatMessage) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	cm.ID = raw["id"].(string)
	cm.ChannelID, _ = uuid.Parse(raw["channelId"].(string))

	segments := raw["segments"].([]interface{})
	for _, segment := range segments {
		segMap := segment.(map[string]interface{})
		segType := MessageSegmentType(segMap["type"].(string))

		var msgSegment MessageSegment
		segData, _ := json.Marshal(segment)

		switch segType {
		case MessageSegmentTypeEmote:
			var emoteSegment MessageSegmentEmote
			if err := json.Unmarshal(segData, &emoteSegment); err != nil {
				return err
			}
			msgSegment = emoteSegment
		case MessageSegmentTypeLink:
			var linkSegment MessageSegmentLink
			if err := json.Unmarshal(segData, &linkSegment); err != nil {
				return err
			}
			msgSegment = linkSegment
		case MessageSegmentTypeMention:
			var mentionSegment MessageSegmentMention
			if err := json.Unmarshal(segData, &mentionSegment); err != nil {
				return err
			}
			msgSegment = mentionSegment
		case MessageSegmentTypeText:
			var textSegment MessageSegmentText
			if err := json.Unmarshal(segData, &textSegment); err != nil {
				return err
			}
			msgSegment = textSegment
		}

		cm.Segments = append(cm.Segments, msgSegment)
	}

	senderData, _ := json.Marshal(raw["sender"])
	var sender User
	if err := json.Unmarshal(senderData, &sender); err != nil {
		return err
	}
	cm.Sender = &sender

	createdAt, _ := time.Parse(time.RFC3339, raw["createdAt"].(string))
	cm.CreatedAt = createdAt

	reactionsData, _ := json.Marshal(raw["reactions"])
	var reactions []ChatMessageReaction
	if err := json.Unmarshal(reactionsData, &reactions); err != nil {
		return err
	}
	cm.Reactions = reactions

	if raw["replyTo"] != nil {
		replyTo, _ := uuid.Parse(raw["replyTo"].(string))
		cm.ReplyTo = &replyTo
	}

	return nil
}
