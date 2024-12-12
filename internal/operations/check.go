package operations

import (
	"context"
	"encoding/json"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var sseChannel = make(chan map[string]map[string]int)

func CountStatusByType(collection *mongo.Collection, status int, docType string) (int64, error) {
	filter := bson.M{
		"status": status,
		"type":   docType,
	}
	count, err := collection.CountDocuments(context.TODO(), filter)
	return count, err
}

func NotifyFrontend(counts map[string]map[string]int) {
	data := map[string]map[string]int{
		"battery": {
			"ok":     counts["battery"]["ok"],
			"medium": counts["battery"]["medium"],
			"bad":    counts["battery"]["bad"],
		},
		"fluid": {
			"ok":     counts["fluid"]["ok"],
			"medium": counts["fluid"]["medium"],
			"bad":    counts["fluid"]["bad"],
		},
		"critical_temp": {
			"ok":  counts["critical_temp"]["ok"],
			"bad": counts["critical_temp"]["bad"],
		},
		"ac_failure": {
			"ok":  counts["ac_failure"]["ok"],
			"bad": counts["ac_failure"]["bad"],
		},
	}

	_, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshalling data to JSON:", err)
		return
	}

	sseChannel <- map[string]map[string]int{
		"battery":       data["battery"],
		"fluid":         data["fluid"],
		"critical_temp": data["critical_temp"],
		"ac_failure":    data["ac_failure"],
	}
}
