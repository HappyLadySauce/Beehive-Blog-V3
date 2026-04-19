# Beehive Blog v2 事件流与 MCP 设计

## 1. 目标

本文件用于定义 v2 第一阶段的：

- 领域事件
- 异步链路
- MCP 对外能力边界

目标是让搜索、摘要、AI 协作和外部智能体接入有统一协议。

## 2. 为什么需要事件流

v2 中有很多操作不适合同步塞进主请求：

- 内容切片
- 搜索索引更新
- 摘要生成
- AI 输出后处理
- 发布后的派生任务

如果全部同步执行，会导致：

- 请求变慢
- 服务耦合变重
- 失败回滚复杂

所以第一阶段就建议使用事件驱动的异步链路。

## 3. 事件设计原则

### 3.1 事件描述“已发生的事实”

例如：

- `content.created`
- `review.approved`

不建议用命令式命名：

- `rebuild-search-now`

### 3.2 事件载荷保持最小必要信息

事件里只放：

- 事件 ID
- 事件类型
- 业务主键
- 时间
- 基础上下文

不要在事件中塞整篇正文快照。

### 3.3 搜索和 AI 都消费主数据事件

不要让内容服务直接调用搜索服务写索引，也不要直接调用 AI 服务生成摘要。

统一采用：

- `content.*` 事件
- `review.*` 事件

驱动派生任务。

## 4. 第一阶段核心事件

## 4.1 identity 事件

- `user.registered`
- `user.logged_in`
- `agent_client.created`
- `agent_client.updated`

## 4.2 content 事件

- `content.created`
- `content.updated`
- `content.deleted`
- `content.status_changed`
- `content.visibility_changed`
- `content.ai_access_changed`
- `content.revision_created`
- `content.revision_restored`
- `content.relation_changed`
- `attachment.created`
- `comment.created`
- `comment.updated`

## 4.3 review 事件

- `review.created`
- `review.approved`
- `review.rejected`

## 4.4 search/index 事件

- `chunk.generated`
- `summary.generated`
- `search.indexed`
- `search.rebuilt`

## 4.5 agent 事件

- `agent.task.created`
- `agent.output.generated`
- `agent.output.submitted`
- `agent.output.accepted`
- `agent.output.rejected`

## 5. 事件载荷建议

建议统一事件结构：

```json
{
  "event_id": "evt_xxx",
  "event_type": "content.updated",
  "occurred_at": "2026-04-15T12:00:00Z",
  "actor_type": "owner",
  "actor_id": "1",
  "resource_type": "content_item",
  "resource_id": "1001",
  "payload": {
    "revision_id": 2003
  }
}
```

## 6. 核心异步链路

## 6.1 内容更新到搜索索引

```text
content-service
  -> publish content.updated / content.revision_created
  -> indexer-worker consume
  -> build chunks
  -> build summaries
  -> update search document
  -> push to Meilisearch / Elasticsearch
  -> publish search.indexed
```

## 6.2 AI 输出到审阅

```text
agent-service
  -> generate output
  -> persist agent_output
  -> publish agent.output.generated
  -> submit review
  -> publish agent.output.submitted
  -> review-service create review_task
```

## 6.3 审阅通过到内容发布

```text
review-service
  -> review.approved
  -> content-service consume
  -> update content status
  -> publish content.status_changed
  -> indexer-worker reindex
```

## 7. 第一阶段消息基础设施建议

如果要控制复杂度，第一阶段可以采用以下策略之一：

### 方案 A：数据库事件表 + worker 轮询

优点：

- 简单
- 依赖少

适合：

- 初期服务数量不多
- 先验证流程

### 方案 B：Redis Stream

优点：

- 成本适中
- 比轮询事件表更自然

适合：

- 已经使用 Redis
- 需要轻量异步队列

### 方案 C：Kafka / NATS

优点：

- 更适合中长期扩展

不足：

- 初期运维复杂度更高

第一阶段建议：

**优先考虑数据库事件表或 Redis Stream，不急着上重消息中间件。**

## 8. MCP 目标

MCP 在 v2 中的作用不是承载全部业务，而是作为外部智能体的标准接入口。

它应负责：

- 暴露 tools
- 暴露 resources
- 协议适配

它不应负责：

- 持久化核心业务数据
- 复杂业务编排

## 9. 第一阶段 MCP 能力边界

第一阶段建议只开放“读优先、写受控”的能力。

### 9.1 Tools

建议第一阶段提供：

- `search_content`
- `read_content`
- `read_related`
- `summarize_content`
- `create_draft`
- `submit_review`

### 9.2 Resources

建议第一阶段提供：

- `content://items/{id}`
- `content://articles/{slug}`
- `content://projects/{slug}`
- `content://experiences/{slug}`
- `content://timeline`
- `content://tags`

### 9.3 Prompts

第一阶段可以后置，不强依赖 MCP prompt。

## 10. MCP Tool 设计建议

## 10.1 search_content

输入：

- `query`
- `types`
- `limit`
- `visibility_scope`

输出：

- 命中的内容列表
- 标题
- 摘要
- slug / id
- 类型
- 引用片段

## 10.2 read_content

输入：

- `id` 或 `slug`

输出：

- 标题
- 正文
- 标签
- 可见性
- 发布时间
- 关联内容

## 10.3 summarize_content

输入：

- `content_id`
- `summary_type`

输出：

- 摘要结果
- 来源内容

## 10.4 create_draft

输入：

- `title`
- `goal`
- `context_content_ids`
- `context_query`

输出：

- `agent_task_id`
- `agent_output_id`
- 草稿文本
- 引用来源

说明：

- 不直接发布
- 默认进入待审链路

## 10.5 submit_review

输入：

- `agent_output_id`

输出：

- `review_task_id`
- 当前状态

## 11. MCP 权限策略

第一阶段 MCP 读取能力必须遵守和 agent 一致的权限规则。

只能读取：

- `published + public + ai_access=allowed`
- `published + member + ai_access=allowed`

默认不能读取：

- `draft`
- `review`
- `private`
- `ai_access=denied`

## 12. MCP 与服务映射

### search_content

调用：

- `search-service`

### read_content / read_related

调用：

- `content-service`
- `search-service`

### summarize_content / create_draft

调用：

- `agent-service`

### submit_review

调用：

- `agent-service`
- `review-service`

## 13. 审计要求

所有 MCP 调用建议至少记录：

- `agent_client_id`
- tool 名称
- 请求时间
- 请求参数摘要
- 返回状态
- 命中的内容 ID 列表

这样后面你才能知道：

- 哪个智能体用了哪些内容
- 是否存在越权尝试
- 哪些内容最常被 AI 消费

## 14. 第一阶段建议的实现顺序

1. 先实现领域事件模型
2. 再落地内容变更 -> 索引更新链路
3. 再实现 agent 输出 -> review 链路
4. 最后封装 MCP tools

## 15. 当前结论

v2 第一阶段的事件与 MCP 设计可以收束为：

**内容、审阅、搜索、AI 之间通过事件解耦；外部智能体通过 MCP 做标准化接入，但权限规则与主系统保持一致。**
