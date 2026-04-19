# Beehive Blog v3 Gateway 设计

## 1. 设计目标

本文件定义 v3 阶段 `gateway` 的定位、边界与演进方向。

v3 的核心前提已经固定：

- 其他业务微服务不提供 HTTP 层，只暴露 RPC
- `gateway` 是唯一业务 HTTP / WebSocket 出口
- 多 `gateway` 实例场景下，通过 `edge` 做用户到最优 `gateway` 的选路

一句话定义：

**`gateway` 是统一业务接入层与流量治理层，负责透传与连接承载，不负责业务编排。**

## 2. 核心问题

在“纯 RPC 微服务”前提下，如果继续把业务编排放进 `gateway`，会出现这些问题：

- `gateway` 会快速膨胀成新的单体边界
- 任意接口调整都可能牵动 `gateway`
- WebSocket、推送、限流与普通业务逻辑容易相互污染
- 多 `gateway` 实例下，接入、选路、承载职责容易混乱

因此 v3 要收口的关键不是“怎么让 gateway 更会编排”，而是：

1. 让 `gateway` 保持薄接入层
2. 让业务编排下沉到领域服务
3. 让多实例路由由 `edge + etcd + redis` 稳定支撑

## 3. 设计原则

### 3.1 统一入口，纯 RPC 后端

所有外部业务流量统一进入 `gateway`。

除 `gateway` 外，其他业务微服务只暴露 RPC，不再提供 HTTP 层。

### 3.2 gateway 只做透传与治理

`gateway` 只负责：

- HTTP -> RPC 转发
- 鉴权与上下文透传
- 限流、超时、熔断、日志、trace
- WebSocket 接入与实时连接管理

不承担业务编排。

### 3.3 业务编排下沉服务

如果一个接口有明确业务主语，则完整视图由对应主领域服务返回。

该服务可以依赖其他服务完成 RPC 编排。

### 3.4 无主语接口单独成服务

如果一个接口没有明确业务主语，不强行塞入已有服务，也不放在 `gateway` 聚合。

直接新开微服务承接，例如：

- `dashboard-service`
- `query-service`
- `composite-service`

### 3.5 路由与承载分层

多实例场景下：

- `edge` 负责用户接入选路
- `gateway` 负责连接承载与消息下发

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
             | dashboard / query / worker  |
             +-----------------------------+
```

## 5. gateway 的职责边界

### 5.1 负责什么

- 对外唯一业务 HTTP / WS 入口
- HTTP 路由匹配与 RPC 转发
- 鉴权与上下文透传
- 限流、超时、熔断、重试
- trace、日志、指标采集
- WebSocket 握手与连接管理
- 用户连接映射与推送下发
- 统一错误包装

### 5.2 不负责什么

- 不保存业务主数据
- 不直接访问数据库
- 不承载身份、内容、搜索等领域真相
- 不承担业务编排
- 不承担“无主语接口”的聚合实现
- 不承担实例发现总线职责

## 6. 接口归属规则

v3 不再使用 `facade`、`aggregate`、`route manifest` 作为推荐方案。

## 6.1 有明确主语的接口

如果接口有明确业务主语：

- HTTP 路由到对应领域服务
- 由该服务返回完整视图
- 需要跨服务数据时，由该服务内部完成 RPC 编排

示例：

- 订单完整详情 -> `order-service`
- 内容完整详情 -> `content-service`
- 用户资料视图 -> `user-service`

## 6.2 无明确主语的接口

如果接口没有明确业务主语：

- 不放在 `gateway`
- 不强行塞进最接近的现有服务
- 直接单独新开服务承接

示例：

- `dashboard-service`
- `query-service`
- `composite-service`

## 6.3 gateway 路由定义方式

第一阶段直接采用：

**代码注册透传路由。**

不引入：

- 复杂路由配置文件
- `route manifest`
- 自动装配单服务接口

结论：

- `gateway` 维护清晰的 HTTP -> RPC 路由注册代码
- 业务复杂度不通过 gateway 层解决，而是通过服务边界解决

## 7. 服务编排原则

### 7.1 主领域服务负责完整视图

若某个接口属于某个明确领域，则该领域服务负责完整响应。

允许：

- `order-service -> product-service`
- `content-service -> user-service`

前提是：

- 编排结果仍属于主领域服务语义

### 7.2 为什么禁止循环依赖

服务间允许单向依赖，但明确禁止循环依赖。

不允许：

- `A -> B -> A`
- `A -> B -> C -> A`

原因：

- 容易引发级联故障
- 领域边界会变模糊
- 变更成本和联调成本显著升高
- 容易出现隐藏递归或重复调用

目标是让服务依赖保持单向图，而不是形成依赖环。

## 8. 多 gateway 路由设计

多实例场景下，v3 增加 `edge` 作为边缘接入层。

### 8.1 `edge` 职责

- 首次接入选路
- 重连优先路由
- 感知 gateway 健康状态与区域信息
- 返回目标 `wsEndpoint`
- 不长期代理 WebSocket 长连接

### 8.2 `gateway` 职责

- 承载 WebSocket 连接
- 管理用户连接与会话
- 维护本地连接状态
- 消费推送事件并下发给本地连接

### 8.3 `etcd` 与 `redis` 分工

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

## 9. 路由状态模型

### 9.1 gateway registry

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

### 9.2 user route

```text
userId -> {
  gatewayId,
  deviceId?,
  expiresAt
}
```

### 9.3 connection route

```text
connId -> gatewayId
```

### 9.4 gateway online index

```text
gatewayId -> set(userId / connId)
```

## 10. edge 选路策略

路由策略固定为：

1. 已有绑定优先
2. 区域 / 延迟优先
3. 健康优先
4. 负载优先

### 10.1 已有绑定优先

如果用户已有有效 `user -> gateway` 映射：

- 优先回原 `gateway`
- 重连尽量保持连接粘性

### 10.2 区域优先

若无现有绑定：

- 优先同 region / zone
- 优先低延迟实例

### 10.3 负载优先

在候选实例中再参考：

- 当前连接数
- 容量水位
- 错误率
- 权重

## 11. 接入模式

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

## 12. 连接与失效回收

### 12.1 TTL 与续约

- `user -> gateway` 映射必须带 TTL
- gateway registry 必须使用 etcd lease
- gateway 需要周期性续约

### 12.2 断连清理

连接关闭时需要清理：

- `conn -> gateway`
- `user -> gateway` 或其引用关系
- 本地在线索引

### 12.3 宕机兜底

当某个 `gateway` 异常宕机时：

- etcd lease 到期后实例自动摘除
- Redis 中旧映射依靠 TTL 自动失效
- 客户端重连后由 `edge` 重新分配

## 13. WebSocket 与 push 模块

建议在 `gateway` 内部拆出独立 realtime 模块。

推荐结构：

```text
services/gateway/
  internal/
    transport/
    auth/
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

## 14. 错误处理

建议保留明确的网关错误码：

- `gateway.route_not_found`
- `gateway.method_not_allowed`
- `gateway.service_unavailable`
- `gateway.upstream_timeout`
- `edge.gateway_not_available`
- `edge.route_assignment_failed`

## 15. 当前结论

v3 的 `gateway` 设计最终收口为：

- 后端服务纯 RPC
- `gateway` 是唯一 HTTP / WS 出口
- `gateway` 只做透传、鉴权、限流、连接承载
- 业务编排优先下沉到主领域服务内部
- 无明确主语的接口直接单独新开服务
- 多实例接入通过 `edge + etcd + redis` 协同完成
- WebSocket 采用“edge 分配、客户端直连 gateway”模式
