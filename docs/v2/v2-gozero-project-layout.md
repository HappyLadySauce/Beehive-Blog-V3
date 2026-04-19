# Beehive Blog v2 go-zero 项目布局设计

## 1. 目标

本文件定义 v2 当前采用的仓库结构、go-zero 服务布局，以及 `goctl` 的实际使用方式。

当前原则已经收口：

- 对外只保留一个 API 服务：`gateway`
- 核心业务服务统一走 RPC
- 目录结构尽量贴近 `goctl` 生成结果
- 不再预留 `domain/` / `repository/`
- 共享代码统一放在 `pkg/`

## 2. 当前仓库结构

```text
beehive-blog/
  api/
    gateway.api
  proto/
    identity.proto
    content.proto
    search.proto
  services/
    gateway/
    identity/
    content/
    search/
    indexer/
  pkg/
    configs/
    constants/
    contracts/
    events/
    libs/
  sql/
    migrations/
    seeds/
  scripts/
    codegen/
    dev/
    db/
  docker/
  deploy/
  docs/
```

## 3. 目录职责

### 3.1 `api/`

只放对外 HTTP 契约。

第一阶段只有：

- `gateway.api`

后续即使增加更多业务能力，也优先继续收敛到 `gateway` 暴露，而不是给每个服务再开一套 HTTP。

### 3.2 `proto/`

放内部 RPC 契约。

当前第一批 proto：

- `identity.proto`
- `content.proto`
- `search.proto`

后续补：

- `review.proto`
- `agent.proto`

### 3.3 `services/`

放服务实现。

当前结构：

```text
services/
  gateway/   # go-zero API 服务
  identity/  # go-zero RPC 服务
  content/   # go-zero RPC 服务
  search/    # go-zero RPC 服务
  indexer/   # 异步 worker
```

### 3.4 `pkg/`

放跨服务共享包。

只允许放：

- 常量
- 通用响应结构
- 事件 topic
- 配置辅助结构
- 真正通用的无业务歧义工具

不要放：

- 内容业务规则
- 数据库写入逻辑
- 服务专属逻辑

## 4. 服务布局约定

## 4.1 `services/gateway`

`gateway` 使用标准 go-zero API 生成结构：

```text
services/gateway/
  gateway.go
  etc/
    gateway-api.yaml
  internal/
    config/
    handler/
    logic/
    svc/
    types/
```

它的职责是：

- 对外唯一 HTTP 入口
- JWT / 鉴权上下文注入
- 调用内部 RPC 服务
- 统一错误码和响应风格

它不负责：

- 直接操作数据库
- 承载核心业务规则

## 4.2 `services/identity`

`identity` 采用 go-zero RPC 风格：

```text
services/identity/
  beehiveblog.identity.go
  identityservice/
    identityservice.go
  pb/
  internal/
    config/
    logic/
    model/
    server/
    svc/
```

职责：

- 注册
- 登录
- token / refresh token
- 当前用户上下文
- agent client 身份

说明：

- 这类结构应尽量由 `goctl rpc protoc` 生成
- 当前仓库里已先放置占位骨架，正式生成仍依赖 `protoc`

## 4.3 `services/content`

结构与 `identity` 一致：

```text
services/content/
  beehiveblog.content.go
  contentservice/
    contentservice.go
  pb/
  internal/
    config/
    logic/
    model/
    server/
    svc/
```

职责：

- 内容主实体
- 内容版本
- 标签
- 关系
- 附件
- 评论

注意：

- 不再额外创建 `domain/` / `repository/`
- 业务逻辑收敛在 `internal/logic`
- 数据访问收敛在 `internal/model`

## 4.4 `services/search`

结构与 `identity` 一致：

```text
services/search/
  beehiveblog.search.go
  searchservice/
    searchservice.go
  pb/
  internal/
    config/
    logic/
    model/
    server/
    svc/
```

职责：

- 关键词检索
- 搜索结果组装
- related content
- 搜索副本读取

搜索引擎客户端如果后续需要封装，优先放在：

- `internal/svc`
- 或独立成 `internal/model/search`

不要提前造一层泛化目录。

## 4.5 `services/indexer`

worker 不强行套 API / RPC 生成结构。

当前建议：

```text
services/indexer/
  indexer.go
  internal/
    config/
    consumer/
    jobs/
    svc/
```

职责：

- 监听内容变更
- 更新索引
- 生成摘要或切片

## 5. `goctl` 使用方式

## 5.1 适合交给 `goctl`

- `gateway.api` -> API 服务骨架
- `proto/*.proto` -> RPC 服务骨架
- model 基础代码生成

## 5.2 当前实际策略

第一阶段按下面方式落地：

1. 手写并固化 `api/` 与 `proto/`
2. 用 `goctl api go` 生成 `services/gateway`
3. 用 `goctl rpc protoc` 生成 `services/identity`、`services/content`、`services/search`
4. 在生成结果上补业务实现

## 5.3 当前阻塞

当前机器上 `goctl` 已安装，但 `protoc` 还未安装。

因此：

- `gateway` 已经可以真实生成
- RPC 服务当前仍是占位骨架
- 要继续生成 RPC，必须先补 `protoc`

## 6. 为什么不再保留 `domain` / `repository`

原因很直接：

- 这不是 go-zero 默认生成风格
- 当前阶段业务还没复杂到值得强行加层
- 先把 RPC 契约、模型、服务边界跑通更重要

当前收口规则：

- `handler` 只处理 HTTP 输入输出
- `logic` 承担用例逻辑
- `svc` 管依赖
- `model` 管数据访问

后续如果某个服务真的出现复杂规则，再局部补抽象，而不是全仓库预埋。

## 7. 当前结论

v2 当前采用的落地结构是：

**`api + proto + services + pkg`，其中 `gateway` 是唯一 API 服务，核心业务服务统一走 RPC，目录尽量贴近 go-zero 生成结果。**
