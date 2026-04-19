# Beehive Blog v2 API DTO 与错误码规范

## 1. 目标

本文件定义 v2 第一阶段 API 的通用响应结构、分页结构、常见 DTO 规则和错误码约定。

目标：

- 统一前后端接口风格
- 统一 go-zero `types` 设计口径
- 避免每个服务各自定义一套错误响应

## 2. 通用响应结构

建议统一为：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

说明：

- `code=0` 表示业务成功
- 非 0 表示业务错误
- HTTP 状态码仍用于表达协议层结果

## 3. 通用错误响应结构

```json
{
  "code": 4001001,
  "message": "invalid request",
  "data": null,
  "trace_id": "xxx"
}
```

建议：

- 错误响应始终带 `trace_id`
- `message` 对外保持克制，不泄露内部实现

## 4. 分页 DTO

## 4.1 分页请求

建议字段：

- `page`
- `pageSize`
- `keyword`
- `sortBy`
- `sortOrder`

默认值建议：

- `page=1`
- `pageSize=20`

限制建议：

- `pageSize <= 100`

## 4.2 分页响应

```json
{
  "list": [],
  "pagination": {
    "page": 1,
    "pageSize": 20,
    "total": 100,
    "totalPages": 5
  }
}
```

## 5. 通用枚举返回建议

统一返回字符串枚举，不返回魔法数字。

例如：

- `status: "published"`
- `visibility: "public"`
- `aiAccess: "allowed"`
- `role: "member"`

## 6. Content DTO 建议

## 6.1 ContentSummary

用于列表与搜索结果。

建议字段：

- `id`
- `type`
- `title`
- `slug`
- `summary`
- `coverImageUrl`
- `status`
- `visibility`
- `aiAccess`
- `publishedAt`
- `tags`

## 6.2 ContentDetail

用于详情页。

建议字段：

- `id`
- `type`
- `title`
- `slug`
- `summary`
- `bodyMarkdown`
- `coverImageUrl`
- `status`
- `visibility`
- `aiAccess`
- `publishedAt`
- `createdAt`
- `updatedAt`
- `tags`
- `relations`
- `attachments`

## 6.3 ContentCreateRequest

建议字段：

- `type`
- `title`
- `slug`
- `summary`
- `bodyMarkdown`
- `coverImageUrl`
- `visibility`
- `aiAccess`
- `tagIds`

类型扩展字段建议通过：

- `profile`

承载，例如：

```json
{
  "type": "project",
  "title": "Beehive v2",
  "profile": {
    "projectStatus": "active",
    "startDate": "2026-04-15"
  }
}
```

## 6.4 ContentUpdateRequest

与创建类似，但字段支持部分更新。

## 7. Revision DTO 建议

## 7.1 RevisionSummary

- `id`
- `revisionNo`
- `editorType`
- `changeSummary`
- `createdAt`

## 7.2 RevisionDetail

- `id`
- `revisionNo`
- `titleSnapshot`
- `summarySnapshot`
- `bodyMarkdown`
- `editorType`
- `createdAt`

## 8. Relation DTO 建议

## 8.1 RelationItem

- `id`
- `relationType`
- `targetContentId`
- `targetType`
- `targetTitle`
- `targetSlug`
- `weight`

## 9. Attachment DTO 建议

## 9.1 AttachmentSummary

- `id`
- `filename`
- `originalFilename`
- `mimeType`
- `fileSize`
- `publicUrl`
- `visibility`
- `aiAccess`
- `createdAt`

## 10. Comment DTO 建议

## 10.1 CommentItem

- `id`
- `contentId`
- `user`
- `parentId`
- `body`
- `status`
- `createdAt`

## 10.2 CommentCreateRequest

- `body`
- `parentId`

## 11. Review DTO 建议

## 11.1 ReviewTaskSummary

- `id`
- `contentId`
- `targetRevisionId`
- `status`
- `submitterType`
- `createdAt`

## 11.2 ReviewDecisionRequest

- `decision`
- `comment`

## 12. Agent DTO 建议

## 12.1 SummarizeRequest

- `contentId`
- `summaryType`

## 12.2 DraftGenerateRequest

- `title`
- `goal`
- `contextContentIds`
- `contextQuery`

## 12.3 AgentOutputDetail

- `id`
- `taskId`
- `outputType`
- `outputText`
- `status`
- `sources`
- `createdAt`

## 13. Search DTO 建议

## 13.1 SearchRequest

- `query`
- `types`
- `tags`
- `page`
- `pageSize`
- `sortBy`

## 13.2 SearchResultItem

- `contentId`
- `type`
- `title`
- `slug`
- `summary`
- `highlight`
- `score`
- `publishedAt`

## 14. 错误码设计原则

建议错误码按模块分段，而不是随意编号。

推荐格式：

- `AA BB CCC`

例如：

- `40 01 001`

解释：

- `AA`：大类
- `BB`：模块
- `CCC`：具体错误

## 15. 错误码分段建议

### 15.1 通用错误 `10xxxxxx`

- `1000001` invalid request
- `1000002` unauthorized
- `1000003` forbidden
- `1000004` resource not found
- `1000005` conflict
- `1000006` too many requests
- `1000007` internal server error

### 15.2 认证模块 `1101xxx`

- `1101001` invalid email
- `1101002` invalid password
- `1101003` user already exists
- `1101004` user not found
- `1101005` password mismatch
- `1101006` token expired
- `1101007` invalid token

### 15.3 内容模块 `1201xxx`

- `1201001` invalid content type
- `1201002` invalid content status
- `1201003` slug already exists
- `1201004` content not found
- `1201005` content not visible
- `1201006` invalid visibility
- `1201007` invalid ai access
- `1201008` invalid content transition

### 15.4 标签与关系模块 `1202xxx`

- `1202001` tag not found
- `1202002` tag already exists
- `1202003` invalid relation type
- `1202004` relation already exists

### 15.5 附件模块 `1203xxx`

- `1203001` file too large
- `1203002` unsupported file type
- `1203003` attachment not found
- `1203004` attachment upload failed

### 15.6 评论模块 `1204xxx`

- `1204001` comment not found
- `1204002` comment forbidden
- `1204003` comment disabled

### 15.7 审阅模块 `1301xxx`

- `1301001` review task not found
- `1301002` review already handled
- `1301003` review decision invalid
- `1301004` review forbidden

### 15.8 搜索模块 `1401xxx`

- `1401001` invalid search query
- `1401002` search backend unavailable
- `1401003` reindex failed

### 15.9 Agent 模块 `1501xxx`

- `1501001` agent task not found
- `1501002` agent output not found
- `1501003` ai access denied
- `1501004` invalid generation context
- `1501005` output submit review failed

## 16. HTTP 状态码建议

建议映射：

- `200` 成功
- `400` 参数错误
- `401` 未认证
- `403` 无权限
- `404` 资源不存在
- `409` 冲突
- `429` 频率限制
- `500` 服务内部错误
- `503` 依赖服务不可用

## 17. go-zero types 建议

在 `internal/types/types.go` 中：

- transport DTO 保持扁平
- 不直接复用数据库 model
- 不直接把内部领域对象泄露到 API

建议每个服务按场景拆 DTO：

- list
- detail
- create
- update
- status change

## 18. 当前结论

v2 第一阶段 API 设计应遵循：

**统一响应、统一分页、统一枚举、统一错误码分段、DTO 面向接口而不是面向数据库。**
