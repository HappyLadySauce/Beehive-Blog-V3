package mq

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestRabbitMQConfigValidation(t *testing.T) {
	t.Parallel()

	if err := (RabbitMQConfig{}).Validate(); err == nil {
		t.Fatalf("expected empty RabbitMQ config to fail")
	}
	cfg := RabbitMQConfig{URL: "amqp://guest:guest@127.0.0.1:5672/", Exchange: "beehive.test.events"}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid RabbitMQ config, got %v", err)
	}
	if got := cfg.WithDefaults(); got.ExchangeType != "topic" || !got.Durable {
		t.Fatalf("unexpected RabbitMQ defaults: %+v", got)
	}
}

func TestRabbitMQPublisherPublishesWhenBrokerConfigured(t *testing.T) {
	url := os.Getenv("BEEHIVE_TEST_RABBITMQ_URL")
	if url == "" {
		t.Skip("BEEHIVE_TEST_RABBITMQ_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	publisher, err := NewRabbitMQPublisher(ctx, RabbitMQConfig{
		URL:      url,
		Exchange: "beehive.test.events",
	})
	if err != nil {
		t.Fatalf("create RabbitMQ publisher failed: %v", err)
	}
	t.Cleanup(func() {
		_ = publisher.Close()
	})

	if err := publisher.Publish(ctx, Message{
		ID:         "evt-test",
		RoutingKey: "content.created",
		Body:       []byte(`{"ok":true}`),
	}); err != nil {
		t.Fatalf("publish RabbitMQ message failed: %v", err)
	}
	if err := publisher.Health(ctx); err != nil {
		t.Fatalf("expected RabbitMQ publisher healthy: %v", err)
	}
}
