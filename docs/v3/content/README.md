# Beehive Blog v3 Content 文档索引

当前 content 文档包括：

- [Content 服务设计](./content-service-design.md)
- [Content 领域模型](./content-domain-model.md)
- [Content 数据库设计](./content-database-design.md)
- [Content API 与 Proto 设计](./content-api-and-proto-design.md)

当前约定：

- `content` 是内容主数据服务，负责内容、版本、标签、关系、附件、评论、状态、可见性和 AI 访问策略。
- `gateway` 只做认证、HTTP 适配和身份透传，不做内容资源级授权。
- `identity` 提供主体身份、角色和账号状态真相，不裁决内容可见性。
- `search`、`indexer`、`agent` 只能消费 content 规则、事件和投影，不成为内容真相源。
- 第一阶段已落地 `items`、`revisions`、`tags`、`content_tags`、`content_relation`、Studio 管理接口和公开读取接口。
- 下一阶段实现优先级为 content events -> search/indexer -> attachments/comments。
