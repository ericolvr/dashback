package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type History struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EquipmentId string             `json:"equipment_id"`
	User        string             `json:"user"`
	Type        string             `json:"type"`
	Message     string             `json:"message"`
	Timestamp   time.Time          `json:"timestamp"`
}
