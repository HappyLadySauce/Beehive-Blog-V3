# Beehive Blog v3 服务契约设计

## 1. 目标

本文件用于定义 v3 阶段各服务的职责边界、调用关系、数据归属、同步与异步链路，以及 `edge`、`gateway`、`WebSocket`、消息队列之间的协作方式。

这份文档的目标不是罗列所有接口细节，而是先把以下问题收口：

- 一个能力应该归哪个服务负责
- 一个接口应该定义在哪一层
- 一个请求应该怎么流转
- 一个异步事件应该由谁发布、谁消费
- `edge`、`gateway`、业务服务、实时模块之间的边界如何划分

## 2. 设计原则

### 2.1 统一入口

所有外部请求统一进入 `gateway` 或先进入 `edge` 再分配到 `gateway`。

### 2.2 后端纯 RPC

除 `gateway` 外，其他业务微服务统一只暴露 RPC。

不再提供业务 HTTP 层。

### 2.3 契约归属 RPC 服务，HTTP 映射归 gateway

业务能力契约优先归 RPC 服务自身。

`gateway` 只维护：

- 对外 HTTP / WS 入口
- HTTP 到 RPC 的映射
- 聚合接口
- 实时接入与推送能力

### 2.4 单一数据归属

每类核心数据只能有一个主归属服务。

其他服务可以缓存、索引、投影，但不能成为真相源。

### 2.5 同步与异步分离

- 查询、写入、鉴权等同步交互走 HTTP -> Gateway -> RPC
- 索引更新、通知推送、AI 后处理等异步任务走事件总线

### 2.6 实时通道独立

`WebSocket` 视为独立接入能力，由 `gateway` 内的 realtime 模块负责。

## 3. 当前服务总览

v3 第一阶段建议收口为以下模块：

- `edge`
- `gateway`
- `identity`
- `content`
- `search`
- `indexer`
- `realtime`

说明：

- `realtime` 在部署上可以先内嵌在 `gateway` 进程中
- 在逻辑边界上，仍然按独立模块建模
- 后续再补 `review`、`agent`、`notification`

## 4. 总体调用关系

```text
Public Web / Studio / External Apps / Agent
                    |
                    v
                   edge
                    |
                    v
                 gateway
          /          |           \
         v           v            v
   identity       content       search
         \           |            /
          \          |           /
           +---------+----------+
                     |
                     v
                 event bus
                     |
         +-----------+-----------+
         |           |           |
         v           v           v
      indexer    realtime    notification
```

规则：

- 普通 HTTP 请求：`Client -> Gateway -> RPC Service`
- WebSocket 接入：`Client -> Edge -> Gateway`
- 异步派生能力统一走 `event bus`
- 推送：`Business/Event -> Realtime/Push -> Gateway -> Client`

## 5. 服务职责边界

## 5.1 `edge`

### 角色

- 边缘接入与路由决策层

### 职责

- 首次接入选路
- 重连优先路由
- 感知 gateway 健康状态
- 选择最优 gateway
- 返回目标 `wsEndpoint`
- 不长期代理 WebSocket 长连接

### 不负责

- 不承载业务 HTTP 接口
- 不维护业务主数据
- 不维持长连接
- 不承担业务消息分发

## 5.2 `gateway`

### 角色

- 统一业务接入层

### 职责

- 对外业务 HTTP / WS 入口
- route manifest 装配 `proxy/facade`
- 承载 `aggregate`
- 管理连接、用户会话、推送与限流
- trace、日志、错误包装

### 不负责

- 不直接访问数据库
- 不承载身份、内容、搜索等领域真相
- 不让所有接口都进入手写逻辑

### 典型接口归属

- `GET /api/v3/auth/me`
- `GET /healthz`
- `GET /readyz`
- `GET /api/v3/ws/connect`
- 聚合 dashboard 接口

## 5.3 `identity`

### 角色

- 认证与身份服务

### 职责

- 注册
- 登录
- token / refresh token
- 当前用户身份解析
- 用户角色与权限基础信息
- Agent Client 身份管理

### 数据归属

- `user`
- `user_session`
- `refresh_token`
- `agent_client`
- `identity_audit`

### 暴露能力

- `Register`
- `Login`
- `RefreshToken`
- `Logout`
- `GetCurrentUser`
- `ValidateAccessToken`

## 5.4 `content`

### 角色

- 内容主数据服务

### 职责

- 内容主实体管理
- 内容版本管理
- 标签与关系管理
- 附件与评论管理
- 内容状态与可见性控制

### 数据归属

- `content_item`
- `content_revision`
- `content_relation`
- `tag`
- `content_tag`
- `attachment`
- `comment`
- 各类 profile 表

### 暴露能力

- `CreateContent`
- `UpdateContent`
- `GetContent`
- `ListContents`
- `GetPublicContent`
- `UpdateContentStatus`
- `UpdateVisibility`
- `UpdateAiAccess`

## 5.5 `search`

### 角色

- 检索与索引查询服务

### 职责

- 公开检索
- Studio 检索
- 搜索结果组装
- 搜索高亮、过滤、相关内容
- 读取索引副本

### 数据归属

- `search_document`
- `content_chunk`
- `content_summary`

### 暴露能力

- `Search`
- `Suggest`
- `Related`
- `GetIndexedDocument`

## 5.6 `indexer`

### 角色

- 异步索引与内容派生 worker

### 职责

- 消费内容变更事件
- 生成搜索文档
- 生成切片与摘要
- 回写检索副本
- 发布索引完成事件

### 消费事件

- `content.created`
- `content.updated`
- `content.deleted`
- `content.status_changed`
- `content.visibility_changed`
- `content.revision_created`

### 发布事件

- `search.indexed`
- `chunk.generated`
- `summary.generated`

## 5.7 `realtime`

### 角色

- 实时连接与消息下发模块

### 部署建议

- 第一阶段可内嵌于 `gateway`
- 逻辑边界上按独立模块设计

### 职责

- WebSocket 握手
- 连接注册与断开清理
- 用户连接映射
- 订阅关系管理
- 心跳保活
- 消息下发
- 广播范围控制

### 不负责

- 不生成业务消息
- 不持久化业务主数据
- 不替代消息队列

## 6. 接口契约与分类规则

## 6.1 三类模型

`gateway` 的对外 HTTP 接口固定分为：

- `proxy`
- `facade`
- `aggregate`

### `proxy`

- 一个 HTTP 接口对应一个 RPC
- 无额外业务编排

### `facade`

- 最终只调用一个服务
- 需要 DTO/分页/上下文适配

### `aggregate`

- 一个 HTTP 接口调用多个 RPC
- 存在明显编排逻辑

## 6.2 新增接口决策规则

新增一个 HTTP 接口时按下面顺序判断：

1. 是否只调用 1 个 RPC
2. 是否不需要跨服务编排
3. 是否不需要网关层业务决策
4. 是否只需要轻量 DTO / 上下文适配

如果成立，则归 `proxy/facade`。  
否则归 `aggregate`。

## 6.3 route manifest

单服务接口通过 `route manifest` 声明：

- `method`
- `path`
- `kind`
- `service`
- `rpc`
- `auth`
- `handler`（仅 `aggregate` 需要）

规则：

- `proxy/facade` 由 gateway 启动时自动装配
- `aggregate` 仅绑定手写 handler

## 7. 用户路由与实例路由状态模型

## 7.1 gateway registry

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

## 7.2 user route

```text
userId -> {
  gatewayId,
  deviceId?,
  expiresAt
}
```

## 7.3 connection route

```text
connId -> gatewayId
```

## 7.4 gateway online index

```text
gatewayId -> set(userId / connId)
```

## 8. etcd 与 redis 分工

### 8.1 `etcd`

负责：

- gateway 实例注册
- 租约续约
- 实例健康状态
- 权重、机房、区域、容量信息

### 8.2 `redis`

负责：

- 高频在线态
- `user -> gateway`
- `conn -> gateway`
- 在线索引集合
- TTL 自动过期

结论：

- 实例发现与健康主存 `etcd`
- 用户路由映射主存 `redis`

## 9. 路由策略顺序

`edge` 的选路策略固定为：

1. 已有绑定优先
2. 区域 / 延迟优先
3. 健康优先
4. 负载优先

说明：

- 有有效绑定时优先回原 gateway
- 无绑定时优先选择同 region / zone 的健康实例
- 候选实例中再比较连接数、容量、错误率、权重

## 10. 连接与失效回收规则

### 10.1 TTL

- `user -> gateway` 映射必须带 TTL
- `conn -> gateway` 映射必须带 TTL

### 10.2 gateway 续约

- gateway registry 使用 etcd lease
- gateway 需要周期性 keepalive

### 10.3 断连清理

连接关闭时需要清理：

- `conn -> gateway`
- `user -> gateway` 或引用关系
- 本地在线索引

### 10.4 宕机兜底

异常宕机时：

- etcd lease 到期后 gateway 自动摘除
- Redis 中旧映射通过 TTL 自动失效
- 客户端重连重新分配

## 11. 请求链路约定

## 11.1 普通 HTTP

```text
Client
  -> Gateway
  -> RPC Service
  -> Gateway
  -> Client
```

## 11.2 WebSocket

```text
Client
  -> Edge
  -> Gateway
```

## 11.3 推送

```text
Business Service / Worker
  -> Event Bus / Push
  -> Gateway
  -> Client
```

## 11.4 内容更新到索引

```text
Client
  -> Gateway
  -> Content RPC
  -> Event Bus
  -> Indexer
  -> Search RPC
```

## 12. 错误处理边界

## 12.1 `edge` 负责

- 目标 gateway 不可分配
- 候选实例为空
- 路由分配失败

## 12.2 `gateway` 负责

- route not found
- method not allowed
- upstream unavailable
- upstream timeout
- unauthorized
- rate limited

## 12.3 RPC 服务负责

- 参数非法
- 资源不存在
- 状态不允许
- 业务规则冲突
- 权限不足

## 13. 安全边界

### 13.1 `edge`

- 只负责接入选路
- 不持有业务状态真相
- 不维持长连接代理

### 13.2 `gateway`

- 入口鉴权
- 握手鉴权
- 连接级限流
- Header 白名单
- trace id 注入

### 13.3 RPC 服务

- 资源级权限判断
- 业务规则级权限校验
- 数据可见性校验

## 14. v3 第一阶段结论

v3 第一阶段服务契约最终收口为：

**一个边缘接入层 `edge`，一个统一业务出口 `gateway`，三个核心 RPC 服务 `identity/content/search`，一个异步索引模块 `indexer`，一个实时模块 `realtime`。**

其中：

- 后端服务纯 RPC
- 单服务接口走 `proxy/facade`
- 聚合接口走 `aggregate`
- 多实例路由由 `edge + etcd + redis` 协同完成
