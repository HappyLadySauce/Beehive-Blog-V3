# Gateway 路由归属说明

## 目标

统一说明 `gateway` 中两类 HTTP 路由的归属方式，避免后续把所有接口继续堆回 `gateway.api`。

## 1. upstream mapping 路由

适用场景：

- 单一 RPC 透传
- 无额外业务编排
- 无额外副作用
- 公共中间件即可满足

当前承接方式：

- 配置文件：`services/gateway/etc/upstreams/*.yaml`
- 运行入口：`services/gateway/gateway.go`
- 契约来源：下游 `proto/*.proto` + `services/gateway/etc/protos/*.protoset`

当前第一批路由：

- `POST /api/v2/auth/register`
- `POST /api/v2/auth/login`
- `POST /api/v2/auth/refresh`
- `GET /api/v2/public/articles`
- `GET /api/v2/search/query`

## 2. custom logic 路由

适用场景：

- 需要鉴权分组
- 需要跨服务编排
- 需要异步副作用
- 需要特殊错误包装或审计

当前承接方式：

- 契约文件：`api/gateway.api`
- 路由文件：`services/gateway/internal/handler/routes.go`
- 业务实现：`services/gateway/internal/logic/gateway/*`

当前典型路由：

- `GET /api/v2/auth/me`
- `/api/v2/studio/*`

## 3. 新增路由判断规则

优先判断是否满足以下条件：

1. 是否只是单 RPC 调用？
2. 是否不需要额外副作用？
3. 是否不需要特殊鉴权链路？
4. 是否不需要跨服务编排？

如果 4 个问题都为“是”，优先走 upstream mapping。  
否则进入 `gateway.api` + custom logic。
