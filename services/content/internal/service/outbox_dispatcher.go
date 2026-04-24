package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/mq"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/repo"
)

var errPermanentOutboxPublish = errors.New("permanent outbox publish error")

type OutboxDispatcherConfig struct {
	DispatchInterval  time.Duration
	BatchSize         int
	MaxAttempts       int
	RetryDelay        time.Duration
	ProcessingTimeout time.Duration
}

// OutboxDispatcher polls pending events and publishes them to MQ.
// OutboxDispatcher 轮询待投递事件并发布到 MQ。
type OutboxDispatcher struct {
	store     *repo.Store
	publisher mq.Publisher
	clock     func() time.Time
	config    OutboxDispatcherConfig
}

// NewOutboxDispatcher builds a dispatcher with production-safe defaults.
// NewOutboxDispatcher 使用生产安全默认值创建 dispatcher。
func NewOutboxDispatcher(store *repo.Store, publisher mq.Publisher, clock func() time.Time, config OutboxDispatcherConfig) *OutboxDispatcher {
	if clock == nil {
		clock = func() time.Time { return time.Now().UTC() }
	}
	if config.DispatchInterval <= 0 {
		config.DispatchInterval = 2 * time.Second
	}
	if config.BatchSize <= 0 {
		config.BatchSize = 50
	}
	if config.MaxAttempts <= 0 {
		config.MaxAttempts = 5
	}
	if config.RetryDelay <= 0 {
		config.RetryDelay = 10 * time.Second
	}
	if config.ProcessingTimeout <= 0 {
		config.ProcessingTimeout = time.Minute
	}
	return &OutboxDispatcher{store: store, publisher: publisher, clock: clock, config: config}
}

// Start runs the dispatcher loop until ctx is canceled.
// Start 运行 dispatcher 循环，直到 ctx 被取消。
func (d *OutboxDispatcher) Start(ctx context.Context) {
	if d == nil {
		return
	}
	go func() {
		ticker := time.NewTicker(d.config.DispatchInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			if _, err := d.DispatchOnce(ctx); err != nil {
				logs.Ctx(ctx).Error("content_outbox_dispatch_failed", err)
			}
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()
}

// DispatchOnce claims one batch and publishes each event at most once.
// DispatchOnce 领取一批事件，并确保每个事件最多被本轮处理一次。
func (d *OutboxDispatcher) DispatchOnce(ctx context.Context) (int, error) {
	if d == nil || d.store == nil || d.publisher == nil {
		return 0, serviceNotInitialized()
	}
	now := d.clock()
	events, err := d.store.Outbox.ClaimDue(ctx, now, d.config.BatchSize, d.config.ProcessingTimeout)
	if err != nil {
		return 0, internalErr(err)
	}
	var markErrors []error
	for _, event := range events {
		if err := d.publishEvent(ctx, event); err != nil {
			nextAttempts := event.Attempts + 1
			nextRetryAt := d.clock().Add(d.config.RetryDelay)
			maxAttempts := 0
			if errors.Is(err, errPermanentOutboxPublish) {
				maxAttempts = d.config.MaxAttempts
			}
			if markErr := d.store.Outbox.MarkPublishFailed(ctx, event.ID, nextAttempts, maxAttempts, nextRetryAt, truncateError(err.Error()), d.clock()); markErr != nil {
				markErrors = append(markErrors, markErr)
				logs.Ctx(ctx).Error(
					"content_outbox_event_mark_failed_publish_failed",
					markErr,
					logs.String("event_id", event.EventID),
					logs.String("event_type", event.EventType),
				)
				continue
			}
			logs.Ctx(ctx).Error(
				"content_outbox_event_publish_failed",
				err,
				logs.String("event_id", event.EventID),
				logs.String("event_type", event.EventType),
			)
			continue
		}
		if err := d.store.Outbox.MarkDone(ctx, event.ID, d.clock()); err != nil {
			markErrors = append(markErrors, err)
			logs.Ctx(ctx).Error(
				"content_outbox_event_mark_done_failed",
				err,
				logs.String("event_id", event.EventID),
				logs.String("event_type", event.EventType),
			)
			continue
		}
	}
	if len(markErrors) > 0 {
		return len(events), internalErr(errors.Join(markErrors...))
	}
	return len(events), nil
}

func (d *OutboxDispatcher) publishEvent(ctx context.Context, event entity.OutboxEvent) error {
	if event.EventType == "" {
		return fmt.Errorf("%w: event type is empty", errPermanentOutboxPublish)
	}
	return d.publisher.Publish(ctx, mq.Message{
		ID:          event.EventID,
		RoutingKey:  event.EventType,
		ContentType: "application/json",
		Body:        []byte(event.PayloadJSON),
		Timestamp:   event.CreatedAt,
		Headers: mq.Headers{
			"event_type":    event.EventType,
			"resource_type": event.ResourceType,
			"resource_id":   strconv.FormatInt(event.ResourceID, 10),
		},
	})
}

func truncateError(value string) string {
	value = strings.TrimSpace(value)
	if len(value) <= 2048 {
		return value
	}
	return value[:2048]
}
