package manifest

import "time"

// PubsubMessage is a container for a pubsub push message
type PubsubMessage struct {
	Attributes  map[string]string `json:"attributes,omitempty"`
	Data        string            `json:"data,omitempty"`
	MessageID   string            `json:"messageId,omitempty"`
	PublishTime time.Time         `json:"publishTime,omitempty"`
	OrderingKey string            `json:"orderingKey,omitempty"`
}
