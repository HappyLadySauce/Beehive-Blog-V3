# Beehive Blog v3 文档索引

当前 v3 文档包括：

- [Gateway 设计](./gateway/gateway-design.md)
- [服务契约设计](./contracts/service-contracts.md)
- [Edge 与 Gateway 路由设计](./gateway/edge-and-gateway-routing-design.md)

当前约定：

- `v3` 是后续正式设计与落地文档目录
- `v2` 仅作为历史方案与参考资料
- 新的架构决策、接口设计、部署方案统一写入 `v3`
- `edge/gateway/routing` 相关决策以 `v3` 文档为唯一准绳
- `gateway-design.md` 的正式口径是“透传型 gateway”
- `service-contracts.md` 的正式口径是“服务内编排优先”
- `edge-and-gateway-routing-design.md` 的正式口径是“边缘选路 + 多 gateway 连接承载”
- `v3` 文档中不再使用 `facade`、`aggregate`、`route manifest` 作为推荐方案
