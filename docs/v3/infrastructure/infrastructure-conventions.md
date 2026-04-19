# Beehive Blog v3 基础设施约定

## 1. 目标

本文件用于统一 `v3` 阶段基础设施选型与使用边界，避免后续在服务发现、消息队列、缓存与公共组件上出现混用。

本文件重点回答：

- `etcd` 到底由谁接入
- MQ 到底选什么
- `redis` 负责什么
- 哪些能力继续使用 go-zero
- 哪些能力改为第三方库

## 2. 当前收口

`v3` 的基础设施边界正式固定为：

- 服务发现与注册：`go-zero + etcd`
- HTTP / RPC 骨架、日志、错误：`go-zero`
- 数据库访问：`GORM`
- Redis 访问：`go-redis`
- 消息队列：`RabbitMQ + 第三方 Go 客户端`
- 公共抽象：`pkg/*`

## 3. etcd 约定

### 3.1 角色

`etcd` 是基础设施注册发现层。

负责：

- RPC 服务注册与发现
- 服务健康状态
- lease / keepalive
- `edge/gateway` 多实例 registry
- 低频实例元数据

不负责：

- 业务事件流
- 异步任务传递
- 高频在线态

### 3.2 接入方式

`etcd` 继续由 go-zero 体系负责接入。

默认使用：

- `zrpc`
- `RpcServerConf`
- `RpcClientConf`
- `Etcd.Hosts`
- `Etcd.Key`

约束：

- 不单独引入第二套 etcd SDK 作为默认方案
- 若 `edge/gateway` 侧确有特殊 lease / registry 需求，可局部补原生 etcd client
- 但仓库默认标准仍是 go-zero + etcd

## 4. RabbitMQ 约定

### 4.1 角色

`RabbitMQ` 是 v3 第一阶段的标准消息队列。

负责：

- 业务事件总线
- 异步消费
- `indexer` 任务投递
- `realtime / push` 异步下发链路
- SSO / 审计等可异步处理任务
- 后续跨服务解耦型消息

不负责：

- 服务注册发现
- 在线态存储
- gateway 实例健康状态

### 4.2 接入方式

MQ 不使用 go-zero 内置方案。

统一采用：

- 第三方 RabbitMQ Go 客户端
- 仓库内统一封装到 `pkg/mq`

约束：

- 服务侧只能依赖 `pkg/mq` 提供的 publisher / consumer 抽象
- 不允许每个服务自行散乱初始化 RabbitMQ 连接
- 不允许把 RabbitMQ 连接逻辑写进 `logic`

### 4.3 建议抽象

`pkg/mq` 第一阶段至少承载：

- `Publisher`
- `Consumer`
- `Message`
- `Headers`
- `Handler`
- `Ack / Nack`
- `Retry / Dead-letter` 基础策略

## 5. redis 约定

`redis` 继续负责高频状态与短期缓存。

负责：

- 高频在线态
- `user -> gateway`
- `conn -> gateway`
- TTL 路由状态
- 短期缓存

不负责：

- RPC 服务发现
- 事件总线默认实现

## 6. 三者边界

统一收口：

- `etcd`：找服务
- `redis`：找连接 / 在线态 / 短期缓存
- `RabbitMQ`：传事件和任务

任何新功能接入前都必须先判断它属于哪一层，不允许混用。

## 7. 代码落位

实现层统一按下面方式组织：

- `services/<svc>/internal/config`
  - 配置结构定义
- `services/<svc>/internal/svc`
  - 注入：
    - GORM
    - Redis
    - RPC client
    - MQ publisher / consumer
- `services/<svc>/internal/logic`
  - 只做业务编排
- `pkg/mq`
  - RabbitMQ 统一封装

默认不新增：

- `pkg/discovery`

因为 `etcd` 继续由 go-zero / `zrpc` 体系承担。

## 8. 配置约定

运行时配置文件继续按服务存放在：

- `services/<svc>/etc/*.yaml`

对需要 MQ 的服务，配置结构建议包括：

- `RabbitMQ`
  - `URL`
  - `Exchange`
  - `Queue`
  - `RoutingKey`
  - `ConsumerTag`
  - `Prefetch`
  - `DeadLetterExchange`
  - `Retry`

对需要 RPC 服务发现的服务，继续使用：

- `RpcServerConf`
- `RpcClientConf`
- `Etcd`

## 9. 当前结论

`v3` 第一阶段基础设施可以收口为：

**`etcd` 继续走 go-zero，`RabbitMQ` 作为标准 MQ 走第三方库，`redis` 继续负责在线态与短期缓存。**
