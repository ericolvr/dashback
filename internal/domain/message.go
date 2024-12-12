package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Uniorg    string             `json:"uniorg"`
	Panel     string             `json:"panel"`
	Node      int                `json:"node"`
	Status    int                `json:"status"`
	Type      string             `json:"type"`
	Monitored bool               `json:"monitored"`
	Timestamp time.Time          `json:"timestamp"`
}
