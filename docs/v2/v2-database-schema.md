# Beehive Blog v2 数据库初版设计

## 1. 目标

本文件给出 v2 第一阶段 PostgreSQL 数据库的初版结构建议。

它不是最终 SQL 文件，但足够作为：

- 建表设计依据
- migration 拆分依据
- service model 划分依据

## 2. 设计原则

- 主数据进入 PostgreSQL
- 搜索索引属于派生数据，可部分镜像到 PostgreSQL
- 公共字段统一
- 能支持版本、权限、AI 审计

## 3. 通用字段建议

大多数表建议包含：

- `id BIGSERIAL PRIMARY KEY`
- `created_at TIMESTAMPTZ NOT NULL DEFAULT now()`
- `updated_at TIMESTAMPTZ NOT NULL DEFAULT now()`

如需软删除，再按需加：

- `deleted_at TIMESTAMPTZ NULL`

## 4. 用户与身份

### users

```sql
id BIGSERIAL PRIMARY KEY
username VARCHAR(64) NOT NULL UNIQUE
nickname VARCHAR(128) NOT NULL
email VARCHAR(255) NOT NULL UNIQUE
password_hash VARCHAR(255) NOT NULL
avatar_url TEXT NULL
bio TEXT NULL
role VARCHAR(32) NOT NULL
status VARCHAR(32) NOT NULL
last_login_at TIMESTAMPTZ NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

### agent_clients

```sql
id BIGSERIAL PRIMARY KEY
name VARCHAR(128) NOT NULL
provider VARCHAR(64) NOT NULL
client_type VARCHAR(64) NOT NULL
status VARCHAR(32) NOT NULL
api_key_hash VARCHAR(255) NOT NULL
description TEXT NULL
last_used_at TIMESTAMPTZ NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

## 5. 内容主实体

### content_items

```sql
id BIGSERIAL PRIMARY KEY
type VARCHAR(32) NOT NULL
title VARCHAR(255) NOT NULL
slug VARCHAR(255) NOT NULL UNIQUE
status VARCHAR(32) NOT NULL
visibility VARCHAR(32) NOT NULL
ai_access VARCHAR(32) NOT NULL
summary TEXT NULL
cover_image_url TEXT NULL
owner_user_id BIGINT NOT NULL REFERENCES users(id)
author_user_id BIGINT NOT NULL REFERENCES users(id)
source_type VARCHAR(32) NOT NULL
current_revision_id BIGINT NULL
comment_enabled BOOLEAN NOT NULL DEFAULT true
is_featured BOOLEAN NOT NULL DEFAULT false
sort_order INT NOT NULL DEFAULT 0
published_at TIMESTAMPTZ NULL
archived_at TIMESTAMPTZ NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

建议索引：

- `UNIQUE (slug)`
- `INDEX (type, status, visibility)`
- `INDEX (published_at DESC)`
- `INDEX (owner_user_id)`
- `INDEX (author_user_id)`

### content_revisions

```sql
id BIGSERIAL PRIMARY KEY
content_id BIGINT NOT NULL REFERENCES content_items(id)
revision_no INT NOT NULL
title_snapshot VARCHAR(255) NOT NULL
summary_snapshot TEXT NULL
body_markdown TEXT NOT NULL
body_json JSONB NULL
editor_type VARCHAR(32) NOT NULL
editor_user_id BIGINT NULL REFERENCES users(id)
editor_agent_client_id BIGINT NULL REFERENCES agent_clients(id)
change_summary TEXT NULL
source_type VARCHAR(32) NOT NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

建议约束：

- `UNIQUE (content_id, revision_no)`

### project_profiles

```sql
content_id BIGINT PRIMARY KEY REFERENCES content_items(id)
project_status VARCHAR(32) NULL
start_date DATE NULL
end_date DATE NULL
tech_stack_json JSONB NULL
repo_url TEXT NULL
demo_url TEXT NULL
role_name VARCHAR(128) NULL
result_summary TEXT NULL
```

### experience_profiles

```sql
content_id BIGINT PRIMARY KEY REFERENCES content_items(id)
start_date DATE NULL
end_date DATE NULL
stage_label VARCHAR(128) NULL
organization_name VARCHAR(255) NULL
location VARCHAR(255) NULL
experience_summary TEXT NULL
```

### timeline_event_profiles

```sql
content_id BIGINT PRIMARY KEY REFERENCES content_items(id)
event_date TIMESTAMPTZ NOT NULL
event_end_date TIMESTAMPTZ NULL
event_type VARCHAR(64) NULL
location VARCHAR(255) NULL
```

### portfolio_profiles

```sql
content_id BIGINT PRIMARY KEY REFERENCES content_items(id)
work_type VARCHAR(64) NULL
preview_url TEXT NULL
external_url TEXT NULL
```

## 6. 标签与关系

### tags

```sql
id BIGSERIAL PRIMARY KEY
name VARCHAR(64) NOT NULL UNIQUE
slug VARCHAR(128) NOT NULL UNIQUE
description TEXT NULL
color VARCHAR(32) NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

### content_tags

```sql
id BIGSERIAL PRIMARY KEY
content_id BIGINT NOT NULL REFERENCES content_items(id)
tag_id BIGINT NOT NULL REFERENCES tags(id)
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

建议约束：

- `UNIQUE (content_id, tag_id)`

### content_relations

```sql
id BIGSERIAL PRIMARY KEY
from_content_id BIGINT NOT NULL REFERENCES content_items(id)
to_content_id BIGINT NOT NULL REFERENCES content_items(id)
relation_type VARCHAR(64) NOT NULL
weight NUMERIC(8,4) NOT NULL DEFAULT 1
sort_order INT NOT NULL DEFAULT 0
metadata_json JSONB NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

## 7. 附件

### attachments

```sql
id BIGSERIAL PRIMARY KEY
filename VARCHAR(255) NOT NULL
original_filename VARCHAR(255) NOT NULL
mime_type VARCHAR(255) NOT NULL
extension VARCHAR(32) NULL
file_size BIGINT NOT NULL
storage_key TEXT NOT NULL
public_url TEXT NULL
checksum VARCHAR(128) NULL
uploaded_by_user_id BIGINT NOT NULL REFERENCES users(id)
visibility VARCHAR(32) NOT NULL
ai_access VARCHAR(32) NOT NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

### content_attachments

```sql
id BIGSERIAL PRIMARY KEY
content_id BIGINT NOT NULL REFERENCES content_items(id)
attachment_id BIGINT NOT NULL REFERENCES attachments(id)
usage_type VARCHAR(32) NOT NULL
sort_order INT NOT NULL DEFAULT 0
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

## 8. 评论

### comments

```sql
id BIGSERIAL PRIMARY KEY
content_id BIGINT NOT NULL REFERENCES content_items(id)
user_id BIGINT NOT NULL REFERENCES users(id)
parent_id BIGINT NULL REFERENCES comments(id)
body TEXT NOT NULL
status VARCHAR(32) NOT NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

## 9. 审阅

### review_tasks

```sql
id BIGSERIAL PRIMARY KEY
content_id BIGINT NOT NULL REFERENCES content_items(id)
target_revision_id BIGINT NOT NULL REFERENCES content_revisions(id)
submitter_type VARCHAR(32) NOT NULL
submitter_user_id BIGINT NULL REFERENCES users(id)
submitter_agent_client_id BIGINT NULL REFERENCES agent_clients(id)
status VARCHAR(32) NOT NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

### review_decisions

```sql
id BIGSERIAL PRIMARY KEY
review_task_id BIGINT NOT NULL REFERENCES review_tasks(id)
reviewer_user_id BIGINT NOT NULL REFERENCES users(id)
decision VARCHAR(32) NOT NULL
comment TEXT NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

## 10. AI 协作

### agent_tasks

```sql
id BIGSERIAL PRIMARY KEY
task_type VARCHAR(64) NOT NULL
requester_user_id BIGINT NOT NULL REFERENCES users(id)
agent_client_id BIGINT NULL REFERENCES agent_clients(id)
target_content_id BIGINT NULL REFERENCES content_items(id)
input_query TEXT NULL
input_options_json JSONB NULL
status VARCHAR(32) NOT NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

### agent_outputs

```sql
id BIGSERIAL PRIMARY KEY
agent_task_id BIGINT NOT NULL REFERENCES agent_tasks(id)
output_type VARCHAR(64) NOT NULL
target_revision_id BIGINT NULL REFERENCES content_revisions(id)
output_text TEXT NOT NULL
metadata_json JSONB NULL
status VARCHAR(32) NOT NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

### agent_output_sources

```sql
id BIGSERIAL PRIMARY KEY
agent_output_id BIGINT NOT NULL REFERENCES agent_outputs(id)
content_id BIGINT NOT NULL REFERENCES content_items(id)
content_chunk_id BIGINT NULL
source_order INT NOT NULL DEFAULT 0
quote_text TEXT NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

## 11. 搜索与索引派生表

### search_documents

```sql
id BIGSERIAL PRIMARY KEY
content_id BIGINT NOT NULL REFERENCES content_items(id)
revision_id BIGINT NOT NULL REFERENCES content_revisions(id)
title TEXT NOT NULL
summary TEXT NULL
body_text TEXT NOT NULL
type VARCHAR(32) NOT NULL
visibility VARCHAR(32) NOT NULL
ai_access VARCHAR(32) NOT NULL
published_at TIMESTAMPTZ NULL
tags_json JSONB NULL
relations_json JSONB NULL
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

### content_chunks

```sql
id BIGSERIAL PRIMARY KEY
content_id BIGINT NOT NULL REFERENCES content_items(id)
revision_id BIGINT NOT NULL REFERENCES content_revisions(id)
chunk_index INT NOT NULL
chunk_type VARCHAR(32) NOT NULL
heading_path TEXT NULL
text TEXT NOT NULL
token_count INT NOT NULL DEFAULT 0
visibility VARCHAR(32) NOT NULL
ai_access VARCHAR(32) NOT NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

### content_summaries

```sql
id BIGSERIAL PRIMARY KEY
content_id BIGINT NOT NULL REFERENCES content_items(id)
revision_id BIGINT NOT NULL REFERENCES content_revisions(id)
summary_type VARCHAR(32) NOT NULL
summary_text TEXT NOT NULL
source_type VARCHAR(32) NOT NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

## 12. 建议的 migration 拆分

建议按以下顺序拆 migration：

1. `001_users_and_agents.sql`
2. `002_content_core.sql`
3. `003_content_profiles.sql`
4. `004_tags_and_relations.sql`
5. `005_attachments.sql`
6. `006_comments.sql`
7. `007_reviews.sql`
8. `008_agent_tasks.sql`
9. `009_search_derivatives.sql`

## 13. 当前结论

v2 第一阶段数据库结构已经可以围绕：

**用户、内容、版本、关系、附件、评论、审阅、AI、搜索派生数据**

这 9 个领域展开。
