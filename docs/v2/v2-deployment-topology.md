# Beehive Blog v2 部署拓扑设计

## 1. 目标

本文件定义当前架构下的部署拓扑。

新的前提是：

- `gateway` 是唯一对外 HTTP 服务
- `identity/content/search` 是内部 RPC 服务
- `indexer` 是内部 worker

## 2. 第一阶段最小组件

建议最小运行集：

- `gateway`
- `identity`
- `content`
- `search`
- `indexer`
- `postgresql`
- `redis`
- `meilisearch`

说明：

- 搜索引擎第一阶段优先 `Meilisearch`
- `review` 和 `agent` 可以后补

## 3. 当前拓扑

```text
Web Public / Studio / External Apps
                |
                v
             gateway (HTTP)
          /      |       \
         v       v        v
 identity (RPC) content (RPC) search (RPC)
                                ^
                                |
                           indexer (worker)

postgresql <-> identity / content / search
redis      <-> gateway / identity / search / indexer
meilisearch <-> search / indexer
```

## 4. 网络暴露规则

对外暴露：

- `gateway`
- 前端站点

仅内网访问：

- `identity`
- `content`
- `search`
- `indexer`
- `postgresql`
- `redis`
- `meilisearch`

## 5. 各服务依赖

### 5.1 `gateway`

依赖：

- `identity`
- `content`
- `search`
- `redis` 可选

### 5.2 `identity`

依赖：

- `postgresql`
- `redis`

### 5.3 `content`

依赖：

- `postgresql`

### 5.4 `search`

依赖：

- `postgresql`
- `redis`
- `meilisearch`

### 5.5 `indexer`

依赖：

- `postgresql`
- `redis` 或数据库事件表
- `meilisearch`

## 6. 本地开发建议

建议模式：

- 基础依赖走 Docker
- Go 服务本机运行

本地至少拉起：

- `postgresql`
- `redis`
- `meilisearch`

服务进程：

- `services/gateway`
- `services/identity`
- `services/content`
- `services/search`
- `services/indexer`

## 7. 配置组织

当前建议：

- `services/gateway/etc/gateway-api.yaml`
- `services/identity/etc/identity.yaml`
- `services/content/etc/content.yaml`
- `services/search/etc/search.yaml`
- `services/indexer/etc/indexer.yaml`

共享环境变量继续收敛到：

- `deploy/local/.env`
- `deploy/staging/.env`
- `deploy/production/.env`

## 8. 当前阻塞

当前 `gateway` 已经能用 `goctl api go` 真实生成。

RPC 服务正式生成还缺：

- `protoc`

因此当前部署层面已经可以先启动：

- 基础依赖
- `gateway`

而 `identity/content/search` 还处于契约和目录已就位、代码待正式生成的状态。

## 9. 当前结论

v2 当前部署模型已经收口为：

**一个对外 HTTP 网关，多个内部 RPC 服务，加一个异步索引 worker。**
