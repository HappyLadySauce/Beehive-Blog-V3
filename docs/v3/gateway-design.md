# Beehive Blog v3 Gateway 设计

## 1. 设计目标

本文件定义 v3 阶段 `gateway` 的定位、边界与演进方向。

v3 的核心前提已经固定：

- 其他业务微服务不提供 HTTP 层，只暴露 RPC
- `gateway` 是唯一业务 HTTP / WebSocket 出口
- 多 `gateway` 实例场景下，需要通过 `edge` 做用户到最优 `gateway` 的选路

一句话定义：

**`gateway` 是统一业务接入层与流量治理层，不是所有业务接口实现的堆放点。**

## 2. 核心问题

如果在“纯 RPC 微服务”前提下，继续把所有外部接口都做成手写 `gateway` 业务逻辑，会出现这些问题：

- `gateway` 会快速膨胀，变成新的单体边界
- 任何单服务接口调整都要同步修改 `gateway`
- 单服务接口与聚合接口混在一起，评审和维护成本很高
- WebSocket、推送、限流、实时路由会与普通 HTTP 业务逻辑相互污染
- 多 `gateway` 实例下，用户路由和连接承载职责不清晰

所以 v3 必须同时解决两件事：

1. `gateway` 内部如何区分单服务接口与聚合接口
2. 多 `gateway` 实例下，用户应该连接哪一台 `gateway`

## 3. 设计原则

### 3.1 统一入口，纯 RPC 后端

所有外部业务流量统一进入 `gateway`。

除 `gateway` 外，其他业务微服务只暴露 RPC，不再提供 HTTP 层。

### 3.2 单服务接口尽量配置化

一个 HTTP 接口如果本质只是调用单个 RPC，就不应该强制写成一段新的 `gateway` 业务逻辑。

优先采用 `route manifest` 驱动装配。

### 3.3 聚合接口少而精

只有跨服务编排、统一上下文、平台聚合和网关专属治理类接口，才进入 `gateway` 手写逻辑。

### 3.4 实时能力独立建模

`WebSocket` 不视为普通 REST 接口的附属能力，而是 `gateway` 内独立 realtime 模块。

### 3.5 路由与承载分层

多实例场景下：

- `edge` 负责用户接入选路
- `gateway` 负责连接承载与推送下发

### 3.6 安全前置

`gateway` 必须承担统一安全责任，包括：

- 认证与上下文透传
- 限流与防刷
- Header 安全控制
- trace id 注入
- WebSocket 握手校验

## 4. 总体架构

```text
Client / Web / Studio / Agent
            |
            v
   +---------------------------+
   |           Edge            |
   | Route / Region / Health   |
   | Gateway Selection         |
   +-------------+-------------+
                 |
                 v
   +---------------------------+
   |         Gateway           |
   | HTTP Ingress / WS Ingress |
   | Auth / RateLimit / Trace  |
   +-----+-------------+-------+
         |             |
         |             +-------------------+
         |                                 |
         v                                 v
 +---------------+  +---------------+  +---------------+
 | identity      |  | content       |  | search        |
 | RPC service   |  | RPC service   |  | RPC service   |
 +-------+-------+  +-------+-------+  +-------+-------+
         |                  |                  |
         +------------------+------------------+
                            |
                            v
                    +---------------+
                    | Event Bus     |
                    | Redis Stream  |
                    | -> NATS/Kafka |
                    +-------+-------+
                            |
                            v
             +--------------+--------------+
             | indexer / review / agent    |
             | notification / async worker |
             +-----------------------------+
```

## 5. gateway 的职责边界

### 5.1 负责什么

- 对外唯一业务 HTTP / WS 入口
- HTTP 路由匹配与流量转发
- route manifest 装配单服务接口
- 承载聚合接口
- WebSocket 握手与连接管理
- 用户连接映射与推送下发
- 限流、超时、熔断、重试
- trace、日志、指标采集
- 统一错误包装

### 5.2 不负责什么

- 不保存业务主数据
- 不直接访问数据库
- 不承载内容、搜索、身份等核心业务规则
- 不让所有接口都进入手写逻辑
- 不承担实例发现总线职责

## 6. 三类接口模型

为避免后续继续把所有接口都堆到 `gateway` 逻辑层，v3 固定采用三类模型：

### 6.1 `proxy`

定义：

- 一个 HTTP 接口对应一个 RPC 方法
- 无额外业务编排
- 无额外副作用
- 只做通用鉴权、参数绑定、错误映射

示例：

- `POST /api/v3/auth/login -> identity.Login`
- `GET /api/v3/public/articles/:slug -> content.GetPublicContent`
- `GET /api/v3/search/query -> search.Search`

### 6.2 `facade`

定义：

- 仍属于单一业务域
- 最终只调用一个服务
- 但 HTTP 形态与 RPC 形态不完全一致
- 需要 DTO 适配、分页标准化、上下文补齐等轻量转换

示例：

- `GET /api/v3/studio/contents`
- `PUT /api/v3/studio/contents/:id/status`

### 6.3 `aggregate`

定义：

- 一个 HTTP 接口调用多个 RPC
- 存在明显编排逻辑
- 可能需要并发查询、结果聚合、降级或组合响应

示例：

- `GET /api/v3/studio/dashboard`
- `GET /api/v3/auth/me`

## 7. 接口判定规则

新增一个 HTTP 接口时，按下面规则判定：

### 7.1 判定为 `proxy / facade`

满足以下条件则归 `proxy/facade`：

- 只调用 1 个 RPC
- 不依赖其他服务结果
- 不需要跨服务编排
- 不需要网关层独立业务决策

其中：

- 完全透传 -> `proxy`
- 需要轻量 DTO/分页/上下文适配 -> `facade`

### 7.2 判定为 `aggregate`

满足任意一条则归 `aggregate`：

- 调用了 2 个及以上服务
- 需要先调 A，再根据结果调 B
- 需要并发汇总多个服务结果
- 需要网关层显式降级、裁剪或聚合缓存

结论：

**只有 `aggregate` 接口才应该进入 `gateway` 手写逻辑。**

## 8. route manifest 驱动方案

为减少以后频繁修改 `gateway` 核心代码，v3 推荐单服务接口采用 `route manifest` 驱动。

### 8.1 manifest 字段模型

```yaml
- method: GET
  path: /api/v3/public/articles/:slug
  kind: proxy
  service: content
  rpc: GetPublicContent
  auth: optional

- method: GET
  path: /api/v3/studio/contents
  kind: facade
  service: content
  rpc: ListContents
  auth: required

- method: GET
  path: /api/v3/studio/dashboard
  kind: aggregate
  handler: StudioDashboard
  auth: required
```

### 8.2 装配规则

- `proxy/facade`：由 `gateway` 启动时根据 manifest 自动注册
- `aggregate`：仅绑定到手写 handler / logic

### 8.3 设计收益

- 单服务接口不必频繁修改 `gateway` 核心逻辑
- 代码评审时能明确看出接口类别
- 便于后续做校验、生成文档与自动装配

## 9. 契约组织方式

不再推荐：

- 业务服务维护自己的 HTTP routes
- `gateway` 只做前缀转发

因为 v3 的后端已经固定为纯 RPC 微服务。

推荐改为：

```text
api/
  gateway.api                  # 仅保留 gateway 自有聚合接口

proto/
  identity.proto
  content.proto
  search.proto

services/
  gateway/
    routes/
      manifest.yaml            # 单服务接口路由清单
```

说明：

- 业务服务维护自己的 RPC 契约
- `gateway` 维护“外部 HTTP -> 内部 RPC”的映射清单
- `aggregate` 接口仍由 `gateway` 逻辑层实现

## 10. 多 gateway 路由设计

多实例场景下，v3 增加 `edge` 作为边缘接入层。

### 10.1 `edge` 职责

- 首次接入选路
- 重连优先路由
- 感知 gateway 健康状态与区域信息
- 返回目标 `wsEndpoint`
- 不长期代理 WebSocket 长连接

### 10.2 `gateway` 职责

- 承载 WebSocket 连接
- 管理用户连接与会话
- 维护本地连接状态
- 消费推送事件并下发给本地连接

### 10.3 `etcd` 与 `redis` 分工

`etcd`：

- gateway 实例注册
- lease / keepalive
- 健康状态
- 权重、区域、可用区、容量信息

`redis`：

- `user -> gateway`
- `conn -> gateway`
- `gateway -> online users`
- 高频在线态映射与 TTL 失效

## 11. 路由状态模型

### 11.1 gateway registry

```text
gatewayId -> {
  host,
  port,
  region,
  zone,
  status,
  weight,
  capacity
}
```

### 11.2 user route

```text
userId -> {
  gatewayId,
  deviceId?,
  expiresAt
}
```

### 11.3 connection route

```text
connId -> gatewayId
```

### 11.4 gateway online index

```text
gatewayId -> set(userId / connId)
```

## 12. edge 选路策略

路由策略固定为：

1. 已有绑定优先
2. 区域 / 延迟优先
3. 健康优先
4. 负载优先

### 12.1 已有绑定优先

如果用户已有有效 `user -> gateway` 映射：

- 优先回原 `gateway`
- 重连尽量保持连接粘性

### 12.2 区域优先

若无现有绑定：

- 优先同 region / zone
- 优先低延迟实例

### 12.3 负载优先

在候选实例中再参考：

- 当前连接数
- 容量水位
- 错误率
- 权重

## 13. 接入模式

v3 第一阶段固定采用：

**`edge` 分配，客户端直连 `gateway`。**

流程：

1. 客户端先请求 `edge`
2. `edge` 根据路由策略选择目标 `gateway`
3. `edge` 返回 `gatewayId` 与 `wsEndpoint`
4. 客户端直连目标 `gateway`

不采用 “edge 长期反向代理 WebSocket” 的原因：

- edge 压力更大
- 长连接代理成本高
- 复杂度和故障面更大

## 14. 连接与失效回收

### 14.1 TTL 与续约

- `user -> gateway` 映射必须带 TTL
- gateway registry 必须使用 etcd lease
- gateway 需要周期性续约

### 14.2 断连清理

连接关闭时需要清理：

- `conn -> gateway`
- `user -> gateway` 或其引用关系
- 本地在线索引

### 14.3 宕机兜底

当某个 `gateway` 异常宕机时：

- etcd lease 到期后实例自动摘除
- Redis 中旧映射依靠 TTL 自动失效
- 客户端重连后由 `edge` 重新分配

## 15. WebSocket 与 push 模块

建议在 `gateway` 内部拆出独立 realtime 模块。

推荐结构：

```text
services/gateway/
  internal/
    proxy/
    facade/
    aggregate/
    realtime/
      handshake/
      connection/
      subscription/
      push/
    ratelimit/
```

`realtime` 模块负责：

- 握手认证
- 连接注册
- 心跳保活
- 用户连接映射
- 消息下发

`push` 模块负责：

- 消费推送事件
- 根据路由状态命中目标 `gateway`
- 下发给本地连接

## 16. 错误处理

建议保留明确的网关错误码：

- `gateway.route_not_found`
- `gateway.method_not_allowed`
- `gateway.service_unavailable`
- `gateway.upstream_timeout`
- `edge.gateway_not_available`
- `edge.route_assignment_failed`

## 17. 当前结论

v3 的 `gateway` 设计最终收口为：

- 后端服务纯 RPC
- `gateway` 是唯一 HTTP / WS 出口
- 单服务接口优先 `proxy / facade`
- 聚合接口才进入 `aggregate`
- `proxy / facade` 通过 `route manifest` 驱动
- 多实例接入通过 `edge + etcd + redis` 协同完成
- WebSocket 采用“edge 分配、客户端直连 gateway”模式
