# Beehive Blog v2 领域模型设计

## 1. 目标

本文件定义 v2 第一阶段的核心领域对象、字段方向、状态机和实体关系。

设计目标：

- 支撑文章、笔记、项目、经历、时间线事件的统一管理
- 支撑版本、发布、搜索、AI 草稿与审阅
- 支撑公开站、Studio、搜索服务、Agent 接入的统一数据基础

## 2. 设计原则

### 2.1 内容统一抽象，具体类型扩展

v2 不再以 `article` 作为唯一核心实体，而是采用统一内容抽象：

- 一个统一内容主表
- 不同内容类型通过 `type` 区分
- 特定字段通过扩展表补充

### 2.2 主数据与检索副本分离

主数据保留在 PostgreSQL。

搜索索引、切片、摘要、向量等属于派生数据，可由异步任务生成。

### 2.3 关系优先，不写死层级

项目与文章、经历与项目、文章与作品等通过关系模型绑定，而不是目录硬编码。

### 2.4 AI 输出必须可审计

AI 草稿、摘要、周报等都必须能追踪：

- 来源内容
- 触发任务
- 使用的上下文
- 审阅结果

## 3. 核心枚举

## 3.1 主体角色 `role`

- `guest`
- `member`
- `owner`
- `agent`

## 3.2 内容类型 `content_type`

- `article`
- `note`
- `project`
- `experience`
- `timeline_event`
- `insight`
- `portfolio`
- `page`

说明：

- `article`：面向读者的正式文章
- `note`：原始、碎片、偏私密或半私密内容
- `project`：项目实体
- `experience`：某一段人生 / 工作 / 成长阶段
- `timeline_event`：时间线节点
- `insight`：反思、洞察、阶段性认知
- `portfolio`：作品集展示项
- `page`：关于页、联系页等静态信息页

## 3.3 内容状态 `content_status`

- `draft`
- `review`
- `published`
- `archived`

## 3.4 可见性 `visibility`

- `public`
- `member`
- `private`

## 3.5 AI 访问策略 `ai_access`

- `allowed`
- `denied`

## 3.6 编辑者类型 `editor_type`

- `human`
- `agent`
- `system`

## 3.7 来源类型 `source_type`

- `manual`
- `import_v1`
- `import_markdown`
- `agent_generated`
- `agent_assisted`
- `system_generated`

## 3.8 关系类型 `relation_type`

第一阶段建议至少支持：

- `belongs_to`
- `related_to`
- `derived_from`
- `references`
- `part_of`
- `depends_on`
- `attached_to`
- `mentioned_in`
- `timeline_of`

## 4. 核心实体总览

第一阶段建议的核心实体如下：

- `user`
- `agent_client`
- `content_item`
- `content_revision`
- `content_relation`
- `tag`
- `content_tag`
- `attachment`
- `content_attachment`
- `comment`
- `review_task`
- `review_decision`
- `agent_task`
- `agent_output`
- `search_document`
- `content_chunk`
- `content_summary`

## 5. 用户与身份模型

## 5.1 user

作用：

- 表示访客注册后的业务用户

建议字段：

- `id`
- `username`
- `nickname`
- `email`
- `password_hash`
- `avatar_url`
- `bio`
- `role`
- `status`
- `created_at`
- `updated_at`
- `last_login_at`

备注：

- 第一阶段虽然开放注册，但后台审核与平台最高权限仅归 `owner`
- `member` 用于评论、登录可见内容消费、后续收藏与订阅

## 5.2 agent_client

作用：

- 记录外部智能体接入主体

第一阶段先统一归为系统授权主体，但仍建议记录具体客户端信息，便于审计。

建议字段：

- `id`
- `name`
- `provider`
- `client_type`
- `status`
- `api_key_hash`
- `description`
- `created_at`
- `updated_at`
- `last_used_at`

说明：

- 权限模型上归类为 `agent`
- 审计层面仍保留具体 client 记录

## 6. 内容主模型

## 6.1 content_item

作用：

- 平台统一内容主实体

建议字段：

- `id`
- `type`
- `title`
- `slug`
- `status`
- `visibility`
- `ai_access`
- `summary`
- `cover_image_url`
- `owner_user_id`
- `author_user_id`
- `source_type`
- `current_revision_id`
- `published_at`
- `archived_at`
- `created_at`
- `updated_at`

扩展建议字段：

- `comment_enabled`
- `is_featured`
- `sort_order`

字段说明：

- `owner_user_id`：内容归属人，第一阶段通常就是你
- `author_user_id`：作者，可扩展为未来多人
- `current_revision_id`：指向当前有效版本

## 6.2 content_revision

作用：

- 保存内容的历史版本

建议字段：

- `id`
- `content_id`
- `revision_no`
- `title_snapshot`
- `summary_snapshot`
- `body_markdown`
- `body_json`
- `editor_type`
- `editor_user_id`
- `editor_agent_client_id`
- `change_summary`
- `source_type`
- `created_at`

说明：

- `body_markdown` 用于正文存储
- `body_json` 可用于未来富文本结构化扩展
- AI 草稿也作为 revision 存储，而不是单独绕开版本体系

## 6.3 content_type 扩展表

统一内容主表之外，第一阶段建议使用扩展表承载类型特有字段。

### `project_profile`

建议字段：

- `content_id`
- `project_status`
- `start_date`
- `end_date`
- `tech_stack_json`
- `repo_url`
- `demo_url`
- `role_name`
- `result_summary`

### `experience_profile`

建议字段：

- `content_id`
- `start_date`
- `end_date`
- `stage_label`
- `organization_name`
- `location`
- `experience_summary`

### `timeline_event_profile`

建议字段：

- `content_id`
- `event_date`
- `event_end_date`
- `event_type`
- `location`

### `portfolio_profile`

建议字段：

- `content_id`
- `work_type`
- `preview_url`
- `external_url`

## 7. 标签模型

## 7.1 tag

建议字段：

- `id`
- `name`
- `slug`
- `description`
- `color`
- `created_at`
- `updated_at`

## 7.2 content_tag

建议字段：

- `id`
- `content_id`
- `tag_id`
- `created_at`

## 8. 关系模型

## 8.1 content_relation

作用：

- 表示任意内容之间的结构化关系

建议字段：

- `id`
- `from_content_id`
- `to_content_id`
- `relation_type`
- `weight`
- `sort_order`
- `metadata_json`
- `created_at`

典型关系例子：

- 项目 `related_to` 文章
- 文章 `derived_from` 笔记
- 时间线事件 `part_of` 经历
- 作品 `belongs_to` 项目
- 反思 `timeline_of` 某段经历

## 9. 附件模型

## 9.1 attachment

作用：

- 管理平台内文件资源

建议字段：

- `id`
- `filename`
- `original_filename`
- `mime_type`
- `extension`
- `file_size`
- `storage_key`
- `public_url`
- `checksum`
- `uploaded_by_user_id`
- `visibility`
- `ai_access`
- `created_at`

第一阶段附件范围需覆盖：

- 图片
- PDF
- Word / Excel
- 音频
- 视频
- 代码文件
- 网页快照

## 9.2 content_attachment

建议字段：

- `id`
- `content_id`
- `attachment_id`
- `usage_type`
- `sort_order`
- `created_at`

`usage_type` 例如：

- `cover`
- `inline`
- `resource`
- `source_material`

## 10. 评论模型

## 10.1 comment

作用：

- 面向公开内容的互动评论

建议字段：

- `id`
- `content_id`
- `user_id`
- `parent_id`
- `body`
- `status`
- `created_at`
- `updated_at`

第一阶段评论策略：

- 默认不做人工审核拦截
- 可接入简单审核接口或风险检查模块

评论状态第一阶段建议保留：

- `visible`
- `hidden`
- `deleted`

## 11. 审阅模型

## 11.1 review_task

作用：

- 表示一次待审流程

建议字段：

- `id`
- `content_id`
- `target_revision_id`
- `submitter_type`
- `submitter_user_id`
- `submitter_agent_client_id`
- `status`
- `created_at`
- `updated_at`

`status` 建议：

- `pending`
- `approved`
- `rejected`

## 11.2 review_decision

作用：

- 记录一次审阅决定

建议字段：

- `id`
- `review_task_id`
- `reviewer_user_id`
- `decision`
- `comment`
- `created_at`

第一阶段默认唯一 reviewer 就是你自己。

## 12. AI 协作模型

## 12.1 agent_task

作用：

- 记录一次 AI 任务请求

建议字段：

- `id`
- `task_type`
- `requester_user_id`
- `agent_client_id`
- `target_content_id`
- `input_query`
- `input_options_json`
- `status`
- `created_at`
- `updated_at`

`task_type` 建议：

- `summarize`
- `draft_generate`
- `weekly_digest`
- `relation_suggest`
- `tag_suggest`

## 12.2 agent_output

作用：

- 记录 AI 输出结果

建议字段：

- `id`
- `agent_task_id`
- `output_type`
- `target_revision_id`
- `output_text`
- `metadata_json`
- `status`
- `created_at`

建议 `status`：

- `generated`
- `submitted`
- `accepted`
- `rejected`

## 12.3 agent_output_source

作用：

- 记录 AI 输出引用了哪些内容和切片

建议字段：

- `id`
- `agent_output_id`
- `content_id`
- `content_chunk_id`
- `source_order`
- `quote_text`

这是 AI 输出可信度的关键审计表。

## 13. 搜索与索引模型

## 13.1 search_document

作用：

- 提供面向搜索引擎的聚合文档

建议字段：

- `id`
- `content_id`
- `revision_id`
- `title`
- `summary`
- `body_text`
- `type`
- `visibility`
- `ai_access`
- `published_at`
- `tags_json`
- `relations_json`
- `updated_at`

备注：

- 如果使用 Meilisearch / Elasticsearch，这张表也可逻辑存在，物理上由同步任务生成索引文档

## 13.2 content_chunk

作用：

- 将长内容切片，服务搜索、摘要、RAG

建议字段：

- `id`
- `content_id`
- `revision_id`
- `chunk_index`
- `chunk_type`
- `heading_path`
- `text`
- `token_count`
- `visibility`
- `ai_access`
- `created_at`

## 13.3 content_summary

作用：

- 存储不同粒度摘要

建议字段：

- `id`
- `content_id`
- `revision_id`
- `summary_type`
- `summary_text`
- `source_type`
- `created_at`

`summary_type` 建议：

- `short`
- `medium`
- `long`
- `key_points`

## 14. 状态机建议

## 14.1 内容状态流转

```text
draft -> review -> published -> archived
draft -> published
review -> draft
published -> draft
published -> archived
archived -> draft
```

说明：

- 人工内容可以直接从 `draft` 发布
- AI 内容一般走 `draft -> review -> published`

## 14.2 可见性策略

第一阶段建议默认规则：

- `draft` 默认 `private`
- `review` 默认 `private`
- `published` 可以是 `public / member / private`
- `archived` 可保持原可见性，但默认不作为主要列表展示

## 14.3 AI 访问规则

规则建议：

- `draft` 默认 `ai_access=denied`
- `review` 默认 `ai_access=denied`
- `published + member` 可配置 `ai_access=allowed`
- `published + private` 默认 `ai_access=denied`
- 私密经历默认 `private + denied`

## 15. 第一阶段最小实体集

如果需要控制实现范围，第一阶段最小可用实体可以先落这些：

- `user`
- `agent_client`
- `content_item`
- `content_revision`
- `content_relation`
- `tag`
- `content_tag`
- `attachment`
- `content_attachment`
- `comment`
- `review_task`
- `review_decision`

搜索和 AI 相关派生实体可以紧随其后补齐：

- `agent_task`
- `agent_output`
- `agent_output_source`
- `content_chunk`
- `content_summary`

## 16. 服务归属建议

### identity-service

- `user`
- `agent_client`

### content-service

- `content_item`
- `content_revision`
- `content_relation`
- `tag`
- `content_tag`
- `attachment`
- `content_attachment`
- `comment`

### review-service

- `review_task`
- `review_decision`

### agent-service

- `agent_task`
- `agent_output`
- `agent_output_source`

### search-service

- `search_document`
- `content_chunk`
- `content_summary`

## 17. 仍待细化的问题

后续进入数据库设计前，还需要继续定这些点：

1. `page` 是否第一阶段就走统一内容模型
2. 评论是否支持通知、提及、楼中楼层级限制
3. 项目与经历的扩展字段是否需要单独更多表
4. 附件是否需要版本
5. 是否要为搜索单独保留“索引状态表”
6. AI 输出是否需要单独模板版本表

## 18. 当前结论

v2 第一阶段的领域中心已经可以明确为：

**以 `content_item + content_revision + content_relation` 为核心，围绕项目、经历、时间线、搜索、评论、附件和 AI 审阅形成统一知识平台模型。**
