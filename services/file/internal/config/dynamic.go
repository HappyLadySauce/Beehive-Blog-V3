package config

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const etcdConfigPrefix = "/configs/file/"

type DynamicConfig struct {
	MaxUploadBytes      int64    `json:"max_upload_bytes"`
	AllowedContentTypes []string `json:"allowed_content_types"`
	PresignTTLSeconds   int      `json:"presign_ttl_seconds"`
}

type ConfigCache struct {
	mu       sync.RWMutex
	cfg      DynamicConfig
	fallback StorageConf
	client   *clientv3.Client
}

func NewConfigCache(client *clientv3.Client, fallback StorageConf) *ConfigCache {
	c := &ConfigCache{
		client:   client,
		fallback: fallback,
	}
	c.cfg = DynamicConfig{
		MaxUploadBytes:      fallback.MaxUploadBytes,
		AllowedContentTypes: fallback.AllowedContentTypes,
		PresignTTLSeconds:   fallback.PresignTTLSeconds,
	}
	return c
}

func (c *ConfigCache) MaxUploadBytes() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.cfg.MaxUploadBytes > 0 {
		return c.cfg.MaxUploadBytes
	}
	return c.fallback.MaxUploadBytes
}

func (c *ConfigCache) AllowedContentTypes() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.cfg.AllowedContentTypes) > 0 {
		return c.cfg.AllowedContentTypes
	}
	return c.fallback.AllowedContentTypes
}

func (c *ConfigCache) PresignTTLSeconds() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.cfg.PresignTTLSeconds > 0 {
		return c.cfg.PresignTTLSeconds
	}
	return c.fallback.PresignTTLSeconds
}

func (c *ConfigCache) Snapshot() DynamicConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cfg
}

// Start loads current values from etcd and begins watching for changes.
// Start 从 etcd 加载当前值并开始监听变更。
func (c *ConfigCache) Start(ctx context.Context) {
	c.loadAll(ctx)
	go c.watch(ctx)
}

func (c *ConfigCache) loadAll(ctx context.Context) {
	resp, err := c.client.Get(ctx, etcdConfigPrefix, clientv3.WithPrefix())
	if err != nil {
		logs.Ctx(ctx).Warn("config_cache_etcd_load_failed", logs.Err(err))
		return
	}
	if len(resp.Kvs) == 0 {
		logs.Ctx(ctx).Info("config_cache_no_etcd_values_using_yaml_fallback")
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, kv := range resp.Kvs {
		c.applyKV(string(kv.Key), string(kv.Value))
	}
	logs.Ctx(ctx).Info("config_cache_loaded", logs.Int("keys", len(resp.Kvs)))
}

func (c *ConfigCache) watch(ctx context.Context) {
	ch := c.client.Watch(ctx, etcdConfigPrefix, clientv3.WithPrefix())
	for wr := range ch {
		if wr.Err() != nil {
			logs.Ctx(ctx).Warn("config_cache_watch_error", logs.Err(wr.Err()))
			continue
		}
		c.mu.Lock()
		for _, ev := range wr.Events {
			c.applyKV(string(ev.Kv.Key), string(ev.Kv.Value))
		}
		c.mu.Unlock()
		logs.Ctx(ctx).Info("config_cache_reloaded", logs.Int("events", len(wr.Events)))
	}
}

func (c *ConfigCache) applyKV(key, value string) {
	switch {
	case strings.HasSuffix(key, "max_upload_bytes"):
		if v, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64); err == nil && v > 0 {
			c.cfg.MaxUploadBytes = v
		}
	case strings.HasSuffix(key, "allowed_content_types"):
		var types []string
		if err := json.Unmarshal([]byte(value), &types); err == nil {
			c.cfg.AllowedContentTypes = types
		}
	case strings.HasSuffix(key, "presign_ttl_seconds"):
		if v, err := strconv.Atoi(strings.TrimSpace(value)); err == nil && v > 0 {
			c.cfg.PresignTTLSeconds = v
		}
	}
}

// Update writes a config value to etcd.
// Update 将配置值写入 etcd。
func (c *ConfigCache) Update(ctx context.Context, key string, value string) error {
	_, err := c.client.Put(ctx, etcdConfigPrefix+key, value)
	return err
}
