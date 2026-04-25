package mq

import (
	"context"
	"errors"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
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

func TestRabbitMQPublisherHealthDoesNotReconnect(t *testing.T) {
	t.Parallel()

	var openCalls atomic.Int32
	publisher := &RabbitMQPublisher{
		cfg:     RabbitMQConfig{Exchange: "beehive.test.events"},
		conn:    &fakeRabbitConnection{closed: true},
		channel: &fakeRabbitChannel{closed: true},
		open: func(context.Context, RabbitMQConfig) (rabbitConnection, rabbitChannel, error) {
			openCalls.Add(1)
			return &fakeRabbitConnection{}, &fakeRabbitChannel{}, nil
		},
	}

	if err := publisher.Health(context.Background()); err == nil {
		t.Fatalf("expected disconnected health check to fail")
	}
	if got := openCalls.Load(); got != 0 {
		t.Fatalf("expected health check to avoid reconnect, got %d reconnects", got)
	}
}

func TestRabbitMQPublisherPublishReconnectsWhenResourcesClosed(t *testing.T) {
	t.Parallel()

	var openCalls atomic.Int32
	channel := &fakeRabbitChannel{}
	publisher := &RabbitMQPublisher{
		cfg:     RabbitMQConfig{Exchange: "beehive.test.events"},
		conn:    &fakeRabbitConnection{closed: true},
		channel: &fakeRabbitChannel{closed: true},
		open: func(context.Context, RabbitMQConfig) (rabbitConnection, rabbitChannel, error) {
			openCalls.Add(1)
			return &fakeRabbitConnection{}, channel, nil
		},
	}

	if err := publisher.Publish(context.Background(), Message{ID: "evt-reconnect", RoutingKey: "content.created"}); err != nil {
		t.Fatalf("publish after reconnect failed: %v", err)
	}
	if got := openCalls.Load(); got != 1 {
		t.Fatalf("expected one reconnect, got %d", got)
	}
	if got := channel.publishCount(); got != 1 {
		t.Fatalf("expected one published message, got %d", got)
	}
}

func TestRabbitMQPublisherPublishContextErrorKeepsHealthyResources(t *testing.T) {
	t.Parallel()

	conn := &fakeRabbitConnection{}
	channel := &fakeRabbitChannel{publishErr: context.DeadlineExceeded}
	publisher := &RabbitMQPublisher{
		cfg:     RabbitMQConfig{Exchange: "beehive.test.events"},
		conn:    conn,
		channel: channel,
	}

	if err := publisher.Publish(context.Background(), Message{ID: "evt-timeout", RoutingKey: "content.created"}); err == nil {
		t.Fatalf("expected publish timeout error")
	}
	if got := conn.closeCount(); got != 0 {
		t.Fatalf("expected healthy connection to stay open, got %d closes", got)
	}
	if got := channel.closeCount(); got != 0 {
		t.Fatalf("expected healthy channel to stay open, got %d closes", got)
	}
	if err := publisher.Health(context.Background()); err != nil {
		t.Fatalf("expected publisher to remain healthy: %v", err)
	}
}

func TestRabbitMQPublisherPublishNonConnectionErrorKeepsHealthyResources(t *testing.T) {
	t.Parallel()

	conn := &fakeRabbitConnection{}
	channel := &fakeRabbitChannel{publishErr: errors.New("publish rejected")}
	publisher := &RabbitMQPublisher{
		cfg:     RabbitMQConfig{Exchange: "beehive.test.events"},
		conn:    conn,
		channel: channel,
	}

	if err := publisher.Publish(context.Background(), Message{ID: "evt-rejected", RoutingKey: "content.created"}); err == nil {
		t.Fatalf("expected publish error")
	}
	if got := conn.closeCount(); got != 0 {
		t.Fatalf("expected connection to stay open for non-connection error, got %d closes", got)
	}
	if got := channel.closeCount(); got != 0 {
		t.Fatalf("expected channel to stay open for non-connection error, got %d closes", got)
	}
}

func TestRabbitMQPublisherPublishErrorClosesResourcesWhenChannelBecameClosed(t *testing.T) {
	t.Parallel()

	conn := &fakeRabbitConnection{}
	channel := &fakeRabbitChannel{publishErr: errors.New("write failed"), closeOnPublish: true}
	publisher := &RabbitMQPublisher{
		cfg:     RabbitMQConfig{Exchange: "beehive.test.events"},
		conn:    conn,
		channel: channel,
	}

	if err := publisher.Publish(context.Background(), Message{ID: "evt-write-failed", RoutingKey: "content.created"}); err == nil {
		t.Fatalf("expected publish error")
	}
	if got := conn.closeCount(); got != 1 {
		t.Fatalf("expected connection to close after channel becomes closed, got %d closes", got)
	}
}

func TestRabbitMQPublisherPublishClosedErrorClosesResources(t *testing.T) {
	t.Parallel()

	conn := &fakeRabbitConnection{}
	channel := &fakeRabbitChannel{publishErr: amqp.ErrClosed}
	publisher := &RabbitMQPublisher{
		cfg:     RabbitMQConfig{Exchange: "beehive.test.events"},
		conn:    conn,
		channel: channel,
	}

	if err := publisher.Publish(context.Background(), Message{ID: "evt-closed", RoutingKey: "content.created"}); err == nil {
		t.Fatalf("expected publish closed error")
	}
	if got := conn.closeCount(); got != 1 {
		t.Fatalf("expected connection to close after closed error, got %d closes", got)
	}
	if got := channel.closeCount(); got != 1 {
		t.Fatalf("expected channel to close after closed error, got %d closes", got)
	}
}

func TestRabbitMQPublisherConcurrentPublishReconnectsOnce(t *testing.T) {
	t.Parallel()

	var openCalls atomic.Int32
	channel := &fakeRabbitChannel{}
	publisher := &RabbitMQPublisher{
		cfg: RabbitMQConfig{Exchange: "beehive.test.events"},
		open: func(context.Context, RabbitMQConfig) (rabbitConnection, rabbitChannel, error) {
			openCalls.Add(1)
			time.Sleep(25 * time.Millisecond)
			return &fakeRabbitConnection{}, channel, nil
		},
	}

	const workers = 8
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := publisher.Publish(context.Background(), Message{ID: "evt-concurrent", RoutingKey: "content.updated"}); err != nil {
				t.Errorf("publish after reconnect failed: %v", err)
			}
		}()
	}
	wg.Wait()

	if got := openCalls.Load(); got != 1 {
		t.Fatalf("expected one reconnect, got %d", got)
	}
	if got := channel.publishCount(); got != workers {
		t.Fatalf("expected %d published messages, got %d", workers, got)
	}
}

type fakeRabbitConnection struct {
	mu     sync.Mutex
	closed bool
	closeN int
}

func (c *fakeRabbitConnection) IsClosed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closed
}

func (c *fakeRabbitConnection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	c.closeN++
	return nil
}

func (c *fakeRabbitConnection) closeCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closeN
}

type fakeRabbitChannel struct {
	mu             sync.Mutex
	closed         bool
	closeN         int
	publishErr     error
	publishN       int
	closeOnPublish bool
}

func (c *fakeRabbitChannel) IsClosed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closed
}

func (c *fakeRabbitChannel) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	c.closeN++
	return nil
}

func (c *fakeRabbitChannel) PublishWithContext(context.Context, string, string, bool, bool, amqp.Publishing) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.publishN++
	if c.closeOnPublish {
		c.closed = true
	}
	return c.publishErr
}

func (c *fakeRabbitChannel) closeCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closeN
}

func (c *fakeRabbitChannel) publishCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.publishN
}
