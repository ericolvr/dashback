package consumers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Alarmtekgit/websocket/internal/domain"
	"github.com/Alarmtekgit/websocket/internal/operations"
	"github.com/Alarmtekgit/websocket/internal/service"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
)

func DatabaseConsumer(messageService service.MessageService, rabbitMQURL string, collection *mongo.Collection) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer channel.Close()

	q, err := channel.QueueDeclare(
		"equipments",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for msg := range msgs {
		var data domain.Message
		err := json.Unmarshal(msg.Body, &data)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		err = messageService.InsertIfNotExists(context.Background(), &data)
		if err != nil {
			log.Printf("Failed to insert message: %v", err)
		}

		go func() {
			statusCounts := map[string]map[string]int{
				"battery": {
					"ok":     0,
					"medium": 0,
					"bad":    0,
				},
				"fluid": {
					"ok":     0,
					"medium": 0,
					"bad":    0,
				},
				"critical_temp": {
					"ok":  0,
					"bad": 0,
				},
				"ac_failure": {
					"ok":  0,
					"bad": 0,
				},
			}

			// fmt.Printf("Received a message: %v\n", data)
			for _, docType := range []string{"battery", "fluid", "critical_temp", "ac_failure"} {
				var statusesToCheck map[int]string

				if docType == "ac_failure" {
					statusesToCheck = map[int]string{0: "ok", 1: "bad"}
				} else if docType == "critical_temp" {
					statusesToCheck = map[int]string{0: "ok", 1: "bad"}
				} else {
					statusesToCheck = map[int]string{0: "ok", 1: "medium", 2: "bad"}
				}

				for status, label := range statusesToCheck {
					count, err := operations.CountStatusByType(collection, status, docType)
					if err != nil {
						log.Println("Error counting documents for", docType, ":", err)
						continue
					}
					statusCounts[docType][label] = int(count)
				}
			}
			operations.NotifyFrontend(statusCounts)
		}()
	}
}

// 			for status, label := range map[int]string{0: "ok", 1: "medium", 2: "bad"} {
// 				for _, docType := range []string{"battery", "fluid"} {
// 					count, err := operations.CountStatusByType(collection, status, docType)
// 					if err != nil {
// 						log.Println("Error counting documents for", docType, ":", err)
// 						continue
// 					}
// 					statusCounts[docType][label] = int(count)
// 				}
// 			}
// 			operations.NotifyFrontend(statusCounts)
// 		}()

// 	}
// }
