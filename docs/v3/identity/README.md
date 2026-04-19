# Beehive Blog v3 Identity 文档索引

本目录用于收口 `v3` 阶段 `identity` 服务的正式设计。

这些文档用于指导以下内容：

- `v3/proto/identity.proto`
- `v3/api/gateway.api` 中 `/api/v3/auth/*`
- `services/identity` 的后续实现
- `gateway` 与 `identity` 的认证协作方式

当前文档包括：

- [Identity 服务设计](./identity-service-design.md)
- [Identity 领域模型](./identity-domain-model.md)
- [Identity API 与 Proto 设计](./identity-api-and-proto-design.md)

当前约定：

- `identity` 是 `v3` 第一阶段的认证与基础账户服务
- `identity` 第一阶段不单独拆分 `user-service`
- 认证主形态为“本地账号 + SSO 并存”
- `gateway` 只做接入、认证前置和上下文透传
- 业务资源授权由业务服务自行判定，不由 `gateway` 或 `identity` 集中裁决

关联文档：

- [服务契约设计](../contracts/service-contracts.md)
- [Gateway 设计](../gateway/gateway-design.md)
- [Edge 与 Gateway 路由设计](../gateway/edge-and-gateway-routing-design.md)
