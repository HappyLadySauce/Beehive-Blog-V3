package mq

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQPublisher publishes messages to a declared RabbitMQ exchange.
// RabbitMQPublisher 将消息发布到已声明的 RabbitMQ exchange。
type RabbitMQPublisher struct {
	cfg     RabbitMQConfig
	conn    *amqp.Connection
	channel *amqp.Channel
	mu      sync.Mutex
	closed  bool
}

// NewRabbitMQPublisher connects to RabbitMQ and declares the target exchange.
// NewRabbitMQPublisher 连接 RabbitMQ 并声明目标 exchange。
func NewRabbitMQPublisher(ctx context.Context, cfg RabbitMQConfig) (*RabbitMQPublisher, error) {
	cfg = cfg.WithDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	conn, ch, err := openRabbitMQResources(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &RabbitMQPublisher{cfg: cfg, conn: conn, channel: ch}, nil
}

func openRabbitMQResources(ctx context.Context, cfg RabbitMQConfig) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := dialRabbitMQ(ctx, cfg)
	if err != nil {
		return nil, nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, nil, fmt.Errorf("open RabbitMQ channel failed: %w", err)
	}
	if err := ch.ExchangeDeclare(
		cfg.Exchange,
		cfg.ExchangeType,
		cfg.Durable,
		false,
		false,
		false,
		nil,
	); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, nil, fmt.Errorf("declare RabbitMQ exchange failed: %w", err)
	}
	return conn, ch, nil
}

func dialRabbitMQ(ctx context.Context, cfg RabbitMQConfig) (*amqp.Connection, error) {
	done := make(chan struct{})
	defer close(done)
	result := make(chan struct {
		conn *amqp.Connection
		err  error
	})
	go func() {
		conn, err := amqp.DialConfig(cfg.URL, amqp.Config{
			Heartbeat: 10 * time.Second,
			Locale:    "en_US",
		})
		res := struct {
			conn *amqp.Connection
			err  error
		}{conn: conn, err: err}
		select {
		case result <- res:
		case <-done:
			if conn != nil {
				_ = conn.Close()
			}
		}
	}()

	timeout := time.NewTimer(time.Duration(cfg.ConnectTimeoutSeconds) * time.Second)
	defer timeout.Stop()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-timeout.C:
		return nil, fmt.Errorf("connect RabbitMQ timed out")
	case res := <-result:
		if res.err != nil {
			return nil, fmt.Errorf("connect RabbitMQ failed: %w", res.err)
		}
		return res.conn, nil
	}
}

// Publish sends a message to the configured exchange using message.RoutingKey.
// Publish 使用 message.RoutingKey 将消息发送到配置的 exchange。
func (p *RabbitMQPublisher) Publish(ctx context.Context, message Message) error {
	if p == nil {
		return fmt.Errorf("RabbitMQ publisher is not initialized")
	}
	if message.RoutingKey == "" {
		return fmt.Errorf("RabbitMQ routing key is required")
	}
	if len(message.Body) == 0 {
		message.Body = []byte("{}")
	}
	if message.ContentType == "" {
		message.ContentType = "application/json"
	}
	if message.Timestamp.IsZero() {
		message.Timestamp = time.Now().UTC()
	}
	headers := amqp.Table{}
	for key, value := range message.Headers {
		headers[key] = value
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	if err := p.ensureOpenLocked(ctx); err != nil {
		return err
	}
	if err := p.channel.PublishWithContext(ctx, p.cfg.Exchange, message.RoutingKey, p.cfg.Mandatory, false, amqp.Publishing{
		MessageId:    message.ID,
		ContentType:  message.ContentType,
		DeliveryMode: amqp.Persistent,
		Timestamp:    message.Timestamp,
		Headers:      headers,
		Body:         message.Body,
	}); err != nil {
		p.closeResourcesLocked()
		return fmt.Errorf("publish RabbitMQ message failed: %w", err)
	}
	return nil
}

// Health verifies the connection and channel are still open.
// Health 校验连接与 channel 仍然打开。
func (p *RabbitMQPublisher) Health(ctx context.Context) error {
	if p == nil {
		return fmt.Errorf("RabbitMQ publisher is not initialized")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.ensureOpenLocked(ctx)
}

// Close closes the channel and connection.
// Close 关闭 channel 与连接。
func (p *RabbitMQPublisher) Close() error {
	if p == nil {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.closed = true
	return p.closeResourcesLocked()
}

func (p *RabbitMQPublisher) ensureOpenLocked(ctx context.Context) error {
	if p.closed {
		return fmt.Errorf("RabbitMQ publisher is closed")
	}
	if p.conn != nil && !p.conn.IsClosed() && p.channel != nil && !p.channel.IsClosed() {
		return nil
	}
	p.closeResourcesLocked()
	conn, ch, err := openRabbitMQResources(ctx, p.cfg)
	if err != nil {
		return fmt.Errorf("reconnect RabbitMQ publisher failed: %w", err)
	}
	p.conn = conn
	p.channel = ch
	return nil
}

func (p *RabbitMQPublisher) closeResourcesLocked() error {
	var err error
	if p.channel != nil && !p.channel.IsClosed() {
		if closeErr := p.channel.Close(); closeErr != nil {
			err = closeErr
		}
	}
	if p.conn != nil && !p.conn.IsClosed() {
		if closeErr := p.conn.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}
	p.channel = nil
	p.conn = nil
	return err
}
