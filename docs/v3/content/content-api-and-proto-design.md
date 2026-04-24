# Beehive Blog v3 Content API 与 Proto 设计

## 1. 目标

本文件记录 content 服务当前 RPC 契约、gateway HTTP 契约，以及后续接口扩展优先级。

当前真相源：

- RPC：`v3/proto/content.proto`
- HTTP：`v3/api/gateway.api`
- Swagger：`v3/api/gateway.yaml`，只由 goctl 生成

relations 第一阶段已经修改 proto、api 并完成生成。

## 2. 当前 RPC 能力

当前 `content.Content` RPC 已开放：

### 2.1 Studio content

- `CreateContent`
- `UpdateContent`
- `GetContent`
- `ListStudioContents`
- `ArchiveContent`

规则：

- 需要可信内部调用方 metadata
- 需要 authenticated user claims
- service 层最终要求 `role=admin`
- 新建内容默认 `draft + private + ai_access=denied`

### 2.2 Revisions

- `ListContentRevisions`
- `GetContentRevision`

规则：

- 需要 `admin`
- revision 只读，不直接更新
- 当前版本由 `content.current_revision_id` 指向

### 2.3 Tags

- `CreateTag`
- `UpdateTag`
- `DeleteTag`
- `ListTags`

规则：

- 需要 `admin`
- tag name 与 slug 唯一
- 已绑定 tag 删除返回 `CodeContentTagInUse`

### 2.4 Relations

- `CreateContentRelation`
- `DeleteContentRelation`
- `ListContentRelations`

规则：

- 需要 `admin`
- 只管理 `from_content_id = content_id` 的出边关系
- 两端 content 必须存在
- 不允许自关联
- 重复关系返回 `CodeContentRelationAlreadyExists`
- `metadata_json` 非空时必须是合法 JSON

### 2.5 Public content

- `ListPublicContents`
- `GetPublicContentBySlug`

规则：

- 不需要登录主体
- 只返回 `status=published AND visibility=public`
- archived 内容默认不可读

### 2.6 Operations

- `Ping`

规则：

- 用于 service readiness / health 适配

## 3. 当前 HTTP 能力

当前 gateway 已开放：

### 3.1 Studio HTTP

前缀：

```text
/api/v3/studio/content
```

接口：

- `GET /items`
- `POST /items`
- `GET /items/:content_id`
- `PUT /items/:content_id`
- `DELETE /items/:content_id`
- `GET /items/:content_id/revisions`
- `GET /items/:content_id/revisions/:revision_id`
- `GET /items/:content_id/relations`
- `POST /items/:content_id/relations`
- `DELETE /items/:content_id/relations/:relation_id`
- `GET /tags`
- `POST /tags`
- `PUT /tags/:tag_id`
- `DELETE /tags/:tag_id`

规则：

- gateway 负责 bearer auth
- gateway 将身份上下文透传给 content RPC
- content service 做资源授权裁决

### 3.2 Public HTTP

前缀：

```text
/api/v3/public/content
```

接口：

- `GET /items`
- `GET /items/:slug`

规则：

- 不要求 bearer auth
- 只读公开已发布内容

## 4. 错误码边界

当前 content 使用 `pkg/errs` 中的 content 错误码。

常见错误：

- `120101`：invalid argument
- `120102`：invalid content type
- `120103`：invalid content status
- `120104`：invalid visibility
- `120105`：invalid ai access
- `120301`：content access forbidden
- `120401`：invalid content transition
- `120501`：content not found
- `120502`：slug already exists
- `120503`：tag not found
- `120504`：tag already exists
- `120505`：revision not found
- `120506`：tag in use
- `120507`：relation not found
- `120508`：relation already exists
- `129901`：content internal error

规则：

- service 层只返回领域错误
- logic 层映射为 gRPC status
- gateway 层映射为 HTTP 响应
- 不向客户端暴露 SQL 原文、底层 cause 或 gRPC 原始文本

## 5. 下一批接口优先级

### 5.1 第一优先级：content events

建议先不对外暴露 HTTP。

实现方向：

- service 写操作发布内部事件
- 事件通过 `pkg/mq` publisher 抽象投递
- indexer/search/realtime 后续消费事件

建议事件：

- `content.created`
- `content.updated`
- `content.archived`
- `content.status_changed`
- `content.visibility_changed`
- `content.ai_access_changed`
- `content.tag_changed`

### 5.2 第二优先级：attachments / comments

attachments 建议先开放 Studio 管理接口。

comments 建议同时考虑公开发表评论与 Studio 管理接口。

评论权限：

- `member` 可以发表评论
- `admin` 可以隐藏、删除、查看所有评论
- `agent` 不参与评论

## 6. 契约变更规则

后续任何契约变更必须先更新 `.proto` 或 `.api` 真相源，再通过 goctl/protoc 生成。
