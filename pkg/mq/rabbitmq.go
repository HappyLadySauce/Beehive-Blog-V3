package mq

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQPublisher publishes messages to a declared RabbitMQ exchange.
// RabbitMQPublisher 将消息发布到已声明的 RabbitMQ exchange。
type RabbitMQPublisher struct {
	cfg         RabbitMQConfig
	conn        rabbitConnection
	channel     rabbitChannel
	open        rabbitConnector
	mu          sync.Mutex
	reconnectMu sync.Mutex
	closed      bool
}

type rabbitConnector func(ctx context.Context, cfg RabbitMQConfig) (rabbitConnection, rabbitChannel, error)

type rabbitConnection interface {
	IsClosed() bool
	Close() error
}

type rabbitChannel interface {
	IsClosed() bool
	Close() error
	PublishWithContext(ctx context.Context, exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error
}

type amqpConnection struct {
	conn *amqp.Connection
}

func (c *amqpConnection) IsClosed() bool {
	return c == nil || c.conn == nil || c.conn.IsClosed()
}

func (c *amqpConnection) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

type amqpChannel struct {
	channel *amqp.Channel
}

func (c *amqpChannel) IsClosed() bool {
	return c == nil || c.channel == nil || c.channel.IsClosed()
}

func (c *amqpChannel) Close() error {
	if c == nil || c.channel == nil {
		return nil
	}
	return c.channel.Close()
}

func (c *amqpChannel) PublishWithContext(ctx context.Context, exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error {
	return c.channel.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg)
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
	return &RabbitMQPublisher{cfg: cfg, conn: conn, channel: ch, open: openRabbitMQResources}, nil
}

func openRabbitMQResources(ctx context.Context, cfg RabbitMQConfig) (rabbitConnection, rabbitChannel, error) {
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
	return &amqpConnection{conn: conn}, &amqpChannel{channel: ch}, nil
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
	if err := p.ensureOpenLocked(ctx); err != nil {
		p.mu.Unlock()
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
		if shouldCloseAfterPublishError(ctx, err, p.conn, p.channel) {
			p.closeResourcesLocked()
		}
		p.mu.Unlock()
		return fmt.Errorf("publish RabbitMQ message failed: %w", err)
	}
	p.mu.Unlock()
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
	if p.closed {
		return fmt.Errorf("RabbitMQ publisher is closed")
	}
	if !p.isOpenLocked() {
		return fmt.Errorf("RabbitMQ publisher is not connected")
	}
	return nil
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
	if p.isOpenLocked() {
		return nil
	}
	p.closeResourcesLocked()
	p.mu.Unlock()

	p.reconnectMu.Lock()
	defer p.reconnectMu.Unlock()

	p.mu.Lock()
	if p.closed {
		return fmt.Errorf("RabbitMQ publisher is closed")
	}
	if p.isOpenLocked() {
		return nil
	}
	p.mu.Unlock()

	conn, ch, err := p.connector()(ctx, p.cfg)
	if err != nil {
		p.mu.Lock()
		return fmt.Errorf("reconnect RabbitMQ publisher failed: %w", err)
	}

	p.mu.Lock()
	if p.closed {
		_ = closeRabbitMQResources(conn, ch)
		return fmt.Errorf("RabbitMQ publisher is closed")
	}
	p.closeResourcesLocked()
	p.conn = conn
	p.channel = ch
	return nil
}

func (p *RabbitMQPublisher) connector() rabbitConnector {
	if p.open != nil {
		return p.open
	}
	return openRabbitMQResources
}

func (p *RabbitMQPublisher) isOpenLocked() bool {
	return p.conn != nil && !p.conn.IsClosed() && p.channel != nil && !p.channel.IsClosed()
}

func (p *RabbitMQPublisher) closeResourcesLocked() error {
	err := closeRabbitMQResources(p.conn, p.channel)
	p.channel = nil
	p.conn = nil
	return err
}

func closeRabbitMQResources(conn rabbitConnection, channel rabbitChannel) error {
	var err error
	if channel != nil && !channel.IsClosed() {
		if closeErr := channel.Close(); closeErr != nil {
			err = closeErr
		}
	}
	if conn != nil && !conn.IsClosed() {
		if closeErr := conn.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}
	return err
}

func shouldCloseAfterPublishError(ctx context.Context, err error, conn rabbitConnection, channel rabbitChannel) bool {
	if ctx.Err() != nil {
		return false
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	if errors.Is(err, amqp.ErrClosed) || errors.Is(err, io.ErrClosedPipe) || errors.Is(err, net.ErrClosed) {
		return true
	}
	return conn == nil || conn.IsClosed() || channel == nil || channel.IsClosed()
}
