package mq

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Headers carries message metadata across services.
// Headers 在服务之间传递消息元数据。
type Headers map[string]string

// Message is the stable project-level MQ publish payload.
// Message 是项目级稳定的 MQ 发布载荷。
type Message struct {
	ID          string
	RoutingKey  string
	ContentType string
	Body        []byte
	Headers     Headers
	Timestamp   time.Time
}

// Publisher publishes messages and exposes lifecycle health.
// Publisher 发布消息并暴露生命周期健康状态。
type Publisher interface {
	Publish(ctx context.Context, message Message) error
	Health(ctx context.Context) error
	Close() error
}

// RabbitMQConfig configures a RabbitMQ publisher.
// RabbitMQConfig 配置 RabbitMQ publisher。
type RabbitMQConfig struct {
	URL                   string
	Exchange              string
	ExchangeType          string
	Durable               bool
	Mandatory             bool
	ConnectTimeoutSeconds int
}

// WithDefaults returns a copy with safe publisher defaults.
// WithDefaults 返回带安全默认值的副本。
func (c RabbitMQConfig) WithDefaults() RabbitMQConfig {
	if strings.TrimSpace(c.ExchangeType) == "" {
		c.ExchangeType = "topic"
	}
	if c.ConnectTimeoutSeconds <= 0 {
		c.ConnectTimeoutSeconds = 5
	}
	c.Durable = true
	return c
}

// Validate checks the minimal connection and exchange settings.
// Validate 校验最小连接与 exchange 配置。
func (c RabbitMQConfig) Validate() error {
	c = c.WithDefaults()
	if strings.TrimSpace(c.URL) == "" {
		return fmt.Errorf("RabbitMQ.URL is required")
	}
	if strings.TrimSpace(c.Exchange) == "" {
		return fmt.Errorf("RabbitMQ.Exchange is required")
	}
	if c.ExchangeType != "topic" && c.ExchangeType != "direct" && c.ExchangeType != "fanout" && c.ExchangeType != "headers" {
		return fmt.Errorf("RabbitMQ.ExchangeType is invalid")
	}
	return nil
}
