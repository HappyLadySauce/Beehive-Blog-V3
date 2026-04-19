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

### 2.3 gateway 只做 HTTP / WS 接入

`gateway` 只维护：

- 对外 HTTP / WS 入口
- 透传路由
- 鉴权、限流、错误包装
- 连接、会话、推送能力

不承担业务编排。

### 2.4 业务编排优先下沉到领域服务

如果接口有明确业务主语，则完整视图由主领域服务返回。

该服务可以依赖其他服务完成 RPC 编排。

### 2.5 无主语接口单独成服务

如果接口没有明确业务主语，不强行塞入已有服务。

直接新开微服务承接，例如：

- `dashboard-service`
- `query-service`
- `composite-service`

### 2.6 单一数据归属

每类核心数据只能有一个主归属服务。

其他服务可以缓存、索引、投影，但不能成为真相源。

### 2.7 实时通道独立

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
- 后续按需要新增 `review`、`agent`、`notification`
- 无主语接口优先单独建服务，而不是回灌到 `gateway`

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
      indexer    realtime   dashboard/query
```

规则：

- 普通 HTTP 请求：`Client -> Gateway -> RPC Service`
- WebSocket 接入：`Client -> Edge -> Gateway`
- 业务编排：`主领域服务 -> 依赖服务`
- 推送：`Business/Event -> Gateway -> Client`

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
- 不承担业务编排

## 5.2 `gateway`

### 角色

- 统一业务接入层

### 职责

- 对外业务 HTTP / WS 入口
- 透传路由
- 鉴权、限流、错误包装
- 管理连接、用户会话、推送与限流
- trace、日志

### 不负责

- 不直接访问数据库
- 不承载身份、内容、搜索等领域真相
- 不承担业务编排
- 不承接无主语接口的聚合实现

### 典型接口归属

- `/api/v3/auth/*`
- `/api/v3/content/*`
- `/api/v3/search/*`
- `/healthz`
- `/ws`

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

## 5.6 `indexer`

### 角色

- 异步索引与内容派生 worker

### 职责

- 消费内容变更事件
- 生成搜索文档
- 生成切片与摘要
- 回写检索副本
- 发布索引完成事件

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

## 5.8 `dashboard/query/composite service`

### 角色

- 无明确业务主语接口的专用服务

### 适用场景

- Dashboard
- 平台级查询
- 跨多个领域但不适合归属任一现有服务的组合视图

### 原则

- 不放在 `gateway`
- 不强行塞进最接近的旧服务
- 直接独立成服务

## 6. 业务编排原则

## 6.1 主领域服务返回完整视图

如果接口有明确业务主语，则由主领域服务负责完整响应。

例如：

- 订单完整信息 -> `order-service`
- 内容完整信息 -> `content-service`
- 用户资料完整视图 -> `user-service`

该服务可以通过 RPC 依赖其他服务补全数据。

## 6.2 为什么允许服务内编排

这样做的好处：

- `gateway` 保持薄接入层
- 业务规则集中在领域内部
- 领域服务能复用完整视图逻辑
- 接入层与业务层解耦更清晰

## 7. 服务依赖约束

## 7.1 允许单向依赖

允许：

- `A -> B`
- `A -> B -> C`

前提是依赖图保持单向。

## 7.2 禁止循环依赖

明确禁止：

- `A -> B -> A`
- `A -> B -> C -> A`

## 7.3 禁止原因

循环依赖会带来：

- 级联故障更严重
- 领域边界模糊
- 变更与联调成本升高
- 隐藏递归或重复调用
- 发布与排障复杂度上升

结论：

**服务间 RPC 调用允许存在，但必须形成单向图。**

## 8. 用户路由与实例路由状态模型

## 8.1 gateway registry

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

## 8.2 user route

```text
userId -> {
  gatewayId,
  deviceId?,
  expiresAt
}
```

## 8.3 connection route

```text
connId -> gatewayId
```

## 8.4 gateway online index

```text
gatewayId -> set(userId / connId)
```

## 9. etcd 与 redis 分工

### 9.1 `etcd`

负责：

- 通过 go-zero `zrpc` 体系承载服务注册与发现
- gateway 实例注册
- 租约续约
- 实例健康状态
- 权重、机房、区域、容量信息

### 9.2 `redis`

负责：

- 高频在线态
- `user -> gateway`
- `conn -> gateway`
- 在线索引集合
- TTL 自动过期

结论：

- 实例发现与健康主存 `etcd`
- 用户路由映射主存 `redis`

### 9.3 `RabbitMQ`

第一阶段标准 MQ 固定为第三方 `RabbitMQ`。

负责：

- 业务事件总线
- 异步消费
- `indexer` 任务投递
- `realtime / push` 异步投递
- 解耦型业务消息

约束：

- 不使用 go-zero 内置 MQ 方案
- 统一通过 `pkg/mq` 封装接入
- 服务侧只依赖 `svc` 注入的 publisher / consumer

统一边界：

- `etcd`：找服务
- `redis`：找连接与在线态
- `RabbitMQ`：传事件和任务

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

## 11.2 业务编排

```text
Client
  -> Gateway
  -> 主领域服务
  -> 依赖服务
  -> 主领域服务
  -> Gateway
  -> Client
```

## 11.3 WebSocket

```text
Client
  -> Edge
  -> Gateway
```

## 11.4 推送

```text
Business Service / Worker
  -> Gateway
  -> Client
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

**一个边缘接入层 `edge`，一个统一业务出口 `gateway`，若干纯 RPC 业务服务，以及一个实时模块 `realtime`。**

其中：

- 后端服务纯 RPC
- `gateway` 只做透传与接入治理
- 业务编排优先下沉到主领域服务
- 无主语接口直接单独建服务
- 多实例路由由 `edge + etcd + redis` 协同完成
