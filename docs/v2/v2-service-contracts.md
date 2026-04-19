# Beehive Blog v2 微服务契约设计

## 1. 目标

本文件定义 v2 当前阶段的服务边界，以及它们在传输层上的关系。

当前架构约束已经固定：

- 外部请求只进 `gateway`
- 内部核心服务统一走 RPC
- 异步派生链路由 `indexer` 处理

## 2. 当前服务总览

第一阶段已确认的服务：

- `gateway`：API 服务
- `identity`：RPC 服务
- `content`：RPC 服务
- `search`：RPC 服务
- `indexer`：worker

第二阶段再补：

- `review`
- `agent`
- `mcp-server`

## 3. 调用关系

```text
Public Web / Studio / External Clients
                |
                v
             gateway
          /     |     \
         v      v      v
    identity  content  search
                        ^
                        |
                     indexer
```

规则：

- 所有 HTTP 请求统一经由 `gateway`
- `gateway` 调用 `identity/content/search` 的 RPC 能力
- `indexer` 不对外暴露 HTTP

## 4. `gateway`

### 传输角色

- 唯一对外 HTTP API

### 职责

- 路由入口
- 鉴权上下文注入
- 请求校验
- RPC 调用编排
- 统一错误包装

### 不负责

- 不直接连数据库
- 不承载内容规则
- 不做索引写入

## 5. `identity`

### 传输角色

- 内部 RPC 服务

### 职责

- 注册
- 登录
- token / refresh token
- 当前用户信息
- agent client 身份

### 数据归属

- `user`
- `agent_client`

### 暴露能力

- `Register`
- `Login`
- `RefreshToken`
- `GetCurrentUser`

## 6. `content`

### 传输角色

- 内部 RPC 服务

### 职责

- 内容主实体
- 内容版本
- 标签
- 关系
- 附件
- 评论
- 内容状态与可见性控制

### 数据归属

- `content_item`
- `content_revision`
- `project_profile`
- `experience_profile`
- `timeline_event_profile`
- `portfolio_profile`
- `tag`
- `content_tag`
- `content_relation`
- `attachment`
- `content_attachment`
- `comment`

### 暴露能力

- `CreateContent`
- `UpdateContent`
- `GetContent`
- `ListContents`
- `UpdateContentStatus`
- `ListPublicArticles`

规则收口：

- slug 校验只放这里
- 状态流转只放这里
- 关系合法性只放这里

## 7. `search`

### 传输角色

- 内部 RPC 服务

### 职责

- 查询公开内容
- 查询工作台内容
- 组装检索结果
- 读取搜索副本

### 数据归属

- `search_document`
- `content_chunk`
- `content_summary`

### 暴露能力

- `Search`
- `Suggest`
- `Related`

## 8. `indexer`

### 传输角色

- 内部异步 worker

### 职责

- 内容变更后更新索引
- 生成搜索副本
- 生成摘要/切片

### 消费事件

- `content.created`
- `content.updated`
- `content.deleted`
- `content.status_changed`
- `content.visibility_changed`
- `content.ai_access_changed`

### 发布事件

- `search.indexed`
- `summary.generated`
- `chunk.generated`

## 9. `gateway` 对外 API 归属

当前对外接口都通过 `api/gateway.api` 统一定义。

第一批由 `gateway` 暴露的路由包括：

- `/api/v2/auth/*`
- `/api/v2/public/articles`
- `/api/v2/studio/contents*`
- `/api/v2/search/query`
- `/api/v2/healthz`

它们分别转发到：

- `auth/*` -> `identity`
- `public/articles` 与 `studio/contents*` -> `content`
- `search/query` -> `search`

## 10. 当前服务落地顺序

建议顺序：

1. `identity`
2. `content`
3. `gateway`
4. `search`
5. `indexer`

后续再补：

6. `review`
7. `agent`

## 11. 当前结论

v2 第一阶段已经收口为：

**一个 API 网关，三个核心 RPC 服务，一个异步索引 worker。**
