package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/cranes-mentoring/obs-contest/purchase-processor/internal/middleware"
	"github.com/cranes-mentoring/obs-contest/purchase-processor/internal/model"
	"github.com/cranes-mentoring/obs-contest/purchase-processor/internal/service"
	"github.com/mitchellh/mapstructure"
)

type ConsumerGroupHandler struct {
	authService service.UserService
}

func NewConsumerGroupHandler(authService service.UserService) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{authService: authService}
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := session.Context()

	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)

		var debMessage model.DebeziumMessage
		err := json.Unmarshal(message.Value, &debMessage)
		if err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			session.MarkMessage(message, "")
			continue
		}

		if debMessage.Payload.After != nil {
			var afterDoc map[string]interface{}

			err := json.Unmarshal([]byte(*debMessage.Payload.After), &afterDoc)
			if err != nil {
				log.Printf("Error unmarshaling 'after' field: %v", err)
			} else {
				log.Printf("payload before as: %v", afterDoc)

				var purchase model.Purchase

				payloadBytes, err := json.Marshal(afterDoc)
				if err != nil {
					log.Printf("Error marshaling 'after' payloadBytes: %v", err)
				}

				err = json.Unmarshal(payloadBytes, &purchase)
				if err != nil {
					log.Printf("Error unmarshaling 'after' payloadBytes to json: %v", err)
				}

				err = mapstructure.Decode(afterDoc, &purchase)
				if err != nil {
					log.Printf("Error decoding 'after' json to document: %v", err)
				} else {
					log.Printf("Processed Purchase: %+v", purchase)
				}

				traceCtx := context.WithValue(ctx, middleware.TraceIDKey, purchase.TraceID)
				log.Printf("Trace ID: %s", purchase.TraceID)

				_, err = h.authService.FindUser(traceCtx, purchase.Username)
				if err != nil {
					log.Printf("Error finding user: %v", err)
				}
			}
		}

		session.MarkMessage(message, "")
	}

	return nil
}
