package provider

import "strings"

// Registry stores all configured identity providers.
// Registry 保存全部已配置的 identity providers。
type Registry struct {
	providers map[string]Provider
}

// NewRegistry builds a provider registry.
// NewRegistry 构建 provider 注册表。
func NewRegistry(providers ...Provider) *Registry {
	items := make(map[string]Provider, len(providers))
	for _, item := range providers {
		if item == nil {
			continue
		}
		items[strings.ToLower(strings.TrimSpace(item.Name()))] = item
	}

	return &Registry{providers: items}
}

// Get returns a provider by normalized name.
// Get 按规范化名称返回 provider。
func (r *Registry) Get(name string) (Provider, bool) {
	if r == nil {
		return nil, false
	}

	item, ok := r.providers[strings.ToLower(strings.TrimSpace(name))]
	return item, ok
}

// GetCallback returns a callback-capable provider by normalized name.
// GetCallback 按规范化名称返回支持 callback 的 provider。
func (r *Registry) GetCallback(name string) (CallbackProvider, bool) {
	item, ok := r.Get(name)
	if !ok {
		return nil, false
	}

	callback, ok := item.(CallbackProvider)
	return callback, ok
}

// Len returns the number of registered providers.
// Len 返回已注册 provider 的数量。
func (r *Registry) Len() int {
	if r == nil {
		return 0
	}

	return len(r.providers)
}
