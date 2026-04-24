# Beehive Blog v3 Content 领域模型

## 1. 目标

本文件定义 content 服务的领域对象、状态、可见性、AI 访问策略和后续扩展模型。

当前第一阶段的领域中心为：

```text
content_item + content_revision + tag + content_tag
```

后续扩展中心为：

```text
content_relation + attachment + content_attachment + comment
```

## 2. 核心枚举

### 2.1 content_type

当前支持：

- `article`
- `note`
- `project`
- `experience`
- `timeline_event`
- `insight`
- `portfolio`
- `page`

### 2.2 content_status

当前支持：

- `draft`
- `review`
- `published`
- `archived`

状态语义：

- `draft`：草稿，只允许 Studio 管理链路读取
- `review`：待审，第一阶段不进入公开消费链路
- `published`：已发布，可以根据 visibility 进入消费链路
- `archived`：已归档，默认不进入公开消费链路

### 2.3 visibility

当前支持：

- `public`
- `member`
- `private`

语义：

- `public`：访客可读，前提是 `status=published`
- `member`：登录用户可读，第一阶段公开 HTTP API 暂未暴露 member 内容读取
- `private`：仅 Studio 管理链路可读

### 2.4 ai_access

当前支持：

- `allowed`
- `denied`

语义：

- `allowed`：允许 agent/search/RAG 在满足 status 和 visibility 条件时读取
- `denied`：禁止 agent/search/RAG 读取

## 3. 当前实体

### 3.1 content_item

作用：

- 平台统一内容主实体
- 保存标题、slug、状态、可见性、AI 访问策略和当前版本指针

当前字段：

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
- `comment_enabled`
- `is_featured`
- `sort_order`
- `published_at`
- `archived_at`
- `created_at`
- `updated_at`

默认值：

- `status=draft`
- `visibility=private`
- `ai_access=denied`
- `source_type=manual`
- `comment_enabled=true`

### 3.2 content_revision

作用：

- 保存内容正文和版本快照
- 支持历史版本、回滚基础、AI 草稿审计扩展

当前字段：

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

约束：

- 同一内容下 `revision_no` 唯一
- `body_json` 非空时必须是合法 JSON
- 当前版本由 `content_item.current_revision_id` 指向

### 3.3 tag

作用：

- content schema 内的标签资源

当前字段：

- `id`
- `name`
- `slug`
- `description`
- `color`
- `created_at`
- `updated_at`

约束：

- `name` 唯一
- `slug` 唯一
- 已绑定内容的 tag 不允许删除

### 3.4 content_tag

作用：

- 内容与 tag 的多对多绑定

当前字段：

- `id`
- `content_id`
- `tag_id`
- `created_at`

约束：

- `(content_id, tag_id)` 唯一
- 删除 content 时级联删除绑定
- 删除 tag 时若存在绑定则被拒绝

## 4. 下一阶段实体

### 4.1 content_relation

作用：

- 表示内容之间的结构化关系
- 支撑知识图谱、相关内容、经历链路、引用网络

建议字段：

- `id`
- `from_content_id`
- `to_content_id`
- `relation_type`
- `weight`
- `sort_order`
- `metadata_json`
- `created_at`
- `updated_at`

建议关系类型：

- `belongs_to`
- `related_to`
- `derived_from`
- `references`
- `part_of`
- `depends_on`
- `timeline_of`

第一阶段约束：

- 不允许自关联
- 两端内容必须存在
- 默认用唯一约束防止重复关系
- 删除 content 时级联删除相关关系

### 4.2 attachment

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

### 4.3 content_attachment

作用：

- 内容与附件的绑定关系

建议字段：

- `id`
- `content_id`
- `attachment_id`
- `usage_type`
- `sort_order`
- `created_at`

建议 `usage_type`：

- `cover`
- `inline`
- `resource`
- `source_material`

### 4.4 comment

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

建议状态：

- `visible`
- `hidden`
- `deleted`

第一阶段策略：

- `member` 可以发表评论
- `admin` 可以隐藏或删除评论
- `agent` 不参与评论

## 5. 状态与权限默认规则

新内容默认：

- `draft`
- `private`
- `ai_access=denied`

公开读取：

- 只返回 `published + public`

后续 agent 读取：

- 只允许 `published + public/member + ai_access=allowed`
- `private`、`draft`、`review`、`archived` 默认拒绝外部 agent 读取
