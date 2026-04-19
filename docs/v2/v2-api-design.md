# Beehive Blog v2 API 设计草案

## 1. 目标

本文件定义 v2 第一阶段 API 的边界与分组方式。

重点目标：

- 为 Public Web 提供内容消费接口
- 为 Studio 提供管理接口
- 为搜索与 AI 提供独立接口
- 为后续 MCP 接入预留清晰边界

## 2. API 分层

第一阶段建议将 API 分为 5 组：

- `auth APIs`
- `public APIs`
- `studio APIs`
- `search APIs`
- `agent APIs`

## 3. 认证接口

建议前缀：

- `/api/v2/auth`

第一阶段建议接口：

- `POST /register`
- `POST /login`
- `POST /refresh`
- `POST /logout`
- `GET /me`

## 4. Public APIs

建议前缀：

- `/api/v2/public`

建议接口：

- `GET /articles`
- `GET /articles/:slug`
- `GET /projects`
- `GET /projects/:slug`
- `GET /experiences`
- `GET /experiences/:slug`
- `GET /timeline`
- `GET /portfolio`
- `GET /pages/:slug`
- `GET /tags`
- `GET /search`
- `GET /content/:id/comments`
- `POST /content/:id/comments`

## 5. Studio APIs

建议前缀：

- `/api/v2/studio`

### 内容管理

- `GET /contents`
- `POST /contents`
- `GET /contents/:id`
- `PUT /contents/:id`
- `PUT /contents/:id/status`
- `PUT /contents/:id/visibility`
- `PUT /contents/:id/ai-access`
- `DELETE /contents/:id`

### 版本管理

- `GET /contents/:id/revisions`
- `GET /contents/:id/revisions/:revisionId`
- `POST /contents/:id/revisions/:revisionId/restore`

### 关系管理

- `GET /contents/:id/relations`
- `POST /contents/:id/relations`
- `DELETE /contents/:id/relations/:relationId`

### 标签管理

- `GET /tags`
- `POST /tags`
- `PUT /tags/:id`
- `DELETE /tags/:id`

### 附件管理

- `GET /attachments`
- `POST /attachments/upload`
- `GET /attachments/:id`
- `DELETE /attachments/:id`

### 审阅管理

- `GET /reviews`
- `POST /reviews/:id/approve`
- `POST /reviews/:id/reject`

### 评论管理

- `GET /comments`
- `PUT /comments/:id/hide`
- `PUT /comments/:id/show`
- `DELETE /comments/:id`

## 6. Search APIs

建议前缀：

- `/api/v2/search`

建议接口：

- `GET /query`
- `GET /suggest`
- `GET /contents/:id/related`
- `POST /rebuild/:contentId`
- `POST /rebuild-all`

当前网关已落地接口（2026-04-16）：

- `GET /api/v2/search/query`（公开搜索，仅返回 `published + public`）
- `GET /api/v2/studio/search/query`（Studio 搜索，owner 作用域）
- `POST /api/v2/studio/search/index/contents/:id`（按内容重建索引）
- `DELETE /api/v2/studio/search/index/contents/:id`（删除索引文档）

## 7. Agent APIs

建议前缀：

- `/api/v2/agent`

建议接口：

- `POST /summarize`
- `POST /drafts/generate`
- `POST /weekly-digest`
- `POST /relations/suggest`
- `POST /tags/suggest`
- `GET /tasks/:id`
- `GET /outputs/:id`
- `POST /outputs/:id/submit-review`

## 8. 第一阶段响应原则

建议统一响应结构：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

错误响应建议至少区分：

- 参数错误
- 未认证
- 无权限
- 内容不可见
- AI 访问被拒绝
- 审阅状态不允许

## 9. 当前建议

这份 API 设计先作为第一版边界草案。

下一步建议继续细化：

- DTO 字段
- 分页与过滤规范
- 权限校验矩阵
- 事件接口与异步任务接口
- MCP tools 对应关系
