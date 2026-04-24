# Beehive Blog v3 Content 数据库设计

## 1. 目标

本文件记录 content 服务数据库设计，当前已落地迁移，以及后续表的规划边界。

当前 content 数据库对象位于：

- `sql/migrations/v3/content/030_v3_content_items.sql`
- `sql/migrations/v3/content/031_v3_content_revisions.sql`
- `sql/migrations/v3/content/032_v3_content_tags.sql`
- `sql/migrations/v3/content/033_v3_content_relations.sql`
- `sql/migrations/v3/content/034_v3_content_outbox_events.sql`

所有 content 主数据统一使用 PostgreSQL `content` schema。

## 2. 当前已落地表

### 2.1 content.items

作用：

- 内容主实体表

关键约束：

- `slug` 全局唯一
- `type` 限定为当前支持的内容类型
- `status` 限定为 `draft/review/published/archived`
- `visibility` 限定为 `public/member/private`
- `ai_access` 限定为 `allowed/denied`
- `owner_user_id` 和 `author_user_id` 引用 `identity.users`
- `current_revision_id` 引用 `content.revisions`

默认值：

- `status='draft'`
- `visibility='private'`
- `ai_access='denied'`
- `source_type='manual'`
- `comment_enabled=true`
- `is_featured=false`
- `sort_order=0`

索引：

- `ux_content_items_slug`
- `idx_content_items_type_status_visibility`
- `idx_content_items_published_at`
- `idx_content_items_owner_user`
- `idx_content_items_author_user`

### 2.2 content.revisions

作用：

- 内容历史版本表
- 保存正文、快照、编辑者与来源

关键约束：

- `content_id` 引用 `content.items(id)`，删除内容时级联删除版本
- `(content_id, revision_no)` 唯一
- `revision_no > 0`
- `editor_type` 限定为 `human/agent/system`
- `source_type` 限定为 `manual/import_v1/import_markdown/agent_generated/agent_assisted`
- `body_json` 使用 JSONB

索引：

- `ux_content_revisions_content_revision_no`
- `idx_content_revisions_content_created`

### 2.3 content.tags

作用：

- content 服务内的标签资源表

关键约束：

- `name` 唯一
- `slug` 唯一

### 2.4 content.content_tags

作用：

- 内容与标签的多对多绑定表

关键约束：

- `content_id` 引用 `content.items(id)`，删除内容时级联删除绑定
- `tag_id` 引用 `content.tags(id)`，使用 `ON DELETE RESTRICT`
- `(content_id, tag_id)` 唯一

设计说明：

- 删除已绑定 tag 应被数据库拒绝。
- service 层仍需要先检查绑定并返回 `CodeContentTagInUse`。
- 这是“业务友好错误 + 数据库兜底”的双保险。

### 2.5 content.relations

作用：

- 保存内容之间的结构化出边关系。

关键约束：

- `from_content_id` 与 `to_content_id` 均引用 `content.items(id)`，删除 content 时级联删除关系。
- `from_content_id <> to_content_id` 禁止自关联。
- `(from_content_id, to_content_id, relation_type)` 唯一。
- `relation_type` 限定为 `belongs_to/related_to/derived_from/references/part_of/depends_on/timeline_of`。
- `metadata_json` 使用 JSONB。

### 2.6 content.outbox_events

作用：

- 保存 content 服务写操作产生的领域事件。
- 作为 RabbitMQ 发布前的可靠投递缓冲，保证业务数据与事件记录在同一数据库事务内提交。

关键约束：

- `event_id` 全局唯一。
- `event_type` 使用事件路由名，例如 `content.created`。
- `payload_json` 使用 JSONB，但不保存正文全文。
- `status` 限定为 `pending/processing/done/failed`。
- `attempts >= 0`。

投递策略：

- dispatcher 使用 `FOR UPDATE SKIP LOCKED` 领取到期事件。
- 发布成功后标记为 `done` 并写入 `published_at`。
- 发布失败后增加 `attempts`、记录 `last_error`，未超过最大次数则回到 `pending` 等待重试。
- 多 dispatcher 并发运行时，同一事件只允许一个 worker 处理。

## 3. 后续迁移规划

后续迁移编号建议从 `035` 开始。

### 3.1 035_v3_content_attachments.sql

规划表：

- `content.attachments`
- `content.content_attachments`

设计边界：

- `attachments` 管文件资源主数据
- `content_attachments` 管内容与附件绑定
- 文件对象存储细节后续由 storage 配置或独立文件服务决定

### 3.2 036_v3_content_comments.sql

规划表：

- `content.comments`

设计边界：

- 评论归 content 服务
- 第一阶段评论状态建议为 `visible/hidden/deleted`
- `member` 可发评论，`admin` 可管理评论

## 4. 当前不创建的表

relations 第一阶段已创建 `content.relations` 迁移，content events 已创建 `content.outbox_events` 迁移。

当前不创建：

- `content.attachments`
- `content.content_attachments`
- `content.comments`
- `search.search_documents`

这些表必须在对应实现任务中按 contract-first 流程补齐。
