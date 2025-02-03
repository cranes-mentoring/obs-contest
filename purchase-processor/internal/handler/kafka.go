package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/cranes-mentoring/obs-contest/purchase-processor/internal/logging"
	"github.com/cranes-mentoring/obs-contest/purchase-processor/internal/model"
	"github.com/cranes-mentoring/obs-contest/purchase-processor/internal/service"
	"go.opentelemetry.io/otel"

	"github.com/IBM/sarama"
	"github.com/mitchellh/mapstructure"
)

// ConsumerGroupHandler handles Kafka messages for the purchase processor.
type ConsumerGroupHandler struct {
	authService service.UserService
}

// NewConsumerGroupHandler creates an instance of ConsumerGroupHandler.
func NewConsumerGroupHandler(authService service.UserService) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{authService: authService}
}

// Setup runs any initialization logic (optional).
func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup runs any cleanup logic (optional).
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim processes messages from a Kafka consumer group claim.
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Use the context provided by the Kafka consumer session.
	ctx := session.Context()

	// Add trace context to the logger for logging spans.
	logging.AddTraceContextToLogger(ctx)

	// Use the existing trace propagation via the context for each message in the claim.
	for message := range claim.Messages() {
		if err := h.handleMessage(ctx, session, message); err != nil {
			log.Printf("Error processing Kafka message: %v", err)
		}
	}

	return nil
}

// handleMessage processes a single Kafka message and ensures proper tracing and context management.
func (h *ConsumerGroupHandler) handleMessage(ctx context.Context, session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage) error {
	log.Printf("Processing Kafka message: topic = %s, timestamp = %v, value = %s",
		message.Topic, message.Timestamp, string(message.Value))

	tracer := otel.Tracer("purchase-processor")
	newCtx, span := tracer.Start(ctx, "handler.handleMessage")
	defer span.End()

	debMessage, err := h.decodeMessage(message.Value)
	if err != nil {
		log.Printf("Failed to decode Debezium message: %v", err)
		session.MarkMessage(message, "")

		return err
	}

	if debMessage.Payload.After != nil {
		if err := h.processPayload(newCtx, debMessage.Payload.After); err != nil {
			log.Printf("Failed to process 'after' payload: %v", err)
		}
	}

	session.MarkMessage(message, "")

	return nil
}

// decodeMessage decodes a Kafka message into a DebeziumMessage structure.
func (h *ConsumerGroupHandler) decodeMessage(value []byte) (*model.DebeziumMessage, error) {
	var debMessage model.DebeziumMessage
	if err := json.Unmarshal(value, &debMessage); err != nil {

		return nil, err
	}

	return &debMessage, nil
}

// processPayload handles the actual processing of the "after" field from a message payload.
func (h *ConsumerGroupHandler) processPayload(ctx context.Context, afterPayload *string) error {
	var afterDoc map[string]interface{}
	if err := json.Unmarshal([]byte(*afterPayload), &afterDoc); err != nil {
		log.Printf("Failed to unmarshal 'after' payload: %v", err)

		return err
	}

	var purchase model.Purchase
	if err := mapstructure.Decode(afterDoc, &purchase); err != nil {
		log.Printf("Failed to map 'after' payload to Purchase struct: %v", err)

		return err
	}

	if _, err := h.authService.FindUser(ctx, purchase.Username); err != nil {
		log.Printf("Failed to find user: %v", err)

		return err
	}

	return nil
}
