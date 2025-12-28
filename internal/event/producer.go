package event

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/savanrajyaguru/ecommerce-go-order-service/config"
	"github.com/savanrajyaguru/ecommerce-go-order-service/pkg/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type EventProducer struct {
	writer *kafka.Writer
}

var Producer *EventProducer

func InitProducer() {
	cfg := config.AppConfig.Kafka
	if len(cfg.Brokers) == 0 {
		log.Println("Kafka brokers not configured, skipping producer initialization")
		return
	}

	w := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Balancer: &kafka.LeastBytes{},
		// Asynchronous writing
		Async: true,
	}

	Producer = &EventProducer{writer: w}
	log.Println("Kafka producer initialized")
}

func (p *EventProducer) PublishEvent(ctx context.Context, eventType string, payload interface{}, key string) error {
	topic, ok := config.AppConfig.Kafka.Topics[eventType]
	if !ok {
		logger.Log.Warn("Topic not found for event type", zap.String("eventType", eventType))
		return nil // Or error, but we fail open usually
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: data,
		Time:  time.Now(),
	}

	// With Async=true, WriteMessages may return nil if the message is buffered.
	// For critical ordering, we might want sync, but prompt asked for non-blocking.
	// The prompt said "Producer MUST support retries and timeouts" and "non-blocking".
	// standard kafka-go writer does retries internally.

	// Create a context with timeout for the write operation (even if async buffer, ensures we don't block forever if full)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		logger.Log.Error("Failed to publish event", zap.String("topic", topic), zap.Error(err))
		return err
	}

	logger.Log.Info("Event published", zap.String("topic", topic), zap.String("key", key))
	return nil
}

func CloseProducer() {
	if Producer != nil && Producer.writer != nil {
		Producer.writer.Close()
	}
}
