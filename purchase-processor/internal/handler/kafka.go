package handler

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
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
				var purchase model.Purchase
				err = mapstructure.Decode(afterDoc, &purchase)
				if err != nil {
					log.Printf("Error decoding 'after' document: %v", err)
				} else {
					log.Printf("Processed Purchase: %+v", purchase)
				}
			}
		}

		session.MarkMessage(message, "")
	}

	return nil
}
