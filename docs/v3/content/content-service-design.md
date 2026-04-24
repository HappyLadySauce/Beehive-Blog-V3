# Beehive Blog v3 Content 服务设计

## 1. 目标

本文件定义 `content` 服务在 v3 第一阶段的职责边界、调用关系、内部层次、权限落点和后续演进顺序。

`content` 的核心目标是成为内容主数据服务：

- 管理统一内容主实体
- 管理内容版本
- 管理标签、内容标签绑定与内容关系
- 后续管理附件、评论
- 执行内容状态、可见性、AI 访问策略
- 为公开站、Studio、search、indexer、agent 提供可信内容来源

## 2. 设计原则

### 2.1 content 是内容真相源

内容主数据只归 `content` 服务。

其他服务可以保存搜索副本、索引文档、AI 上下文、展示缓存，但不能成为内容最终真相源。

### 2.2 gateway 只做接入

`gateway` 负责：

- HTTP 请求绑定
- access token 校验
- 身份上下文透传
- RPC 调用适配
- HTTP 错误包装

`gateway` 不负责：

- 内容资源级授权
- 内容状态流转规则
- 标签、版本、关系等领域规则
- 直接访问 content 数据库

### 2.3 授权下沉到 content

`content` 按 [第一阶段权限模型](../permission-model.md) 做最终授权裁决。

第一阶段规则：

- Studio 内容管理只允许 `admin`
- `member` 和 `guest` 不允许访问 Studio 内容管理
- 公开读取只返回 `published + public`
- `archived` 默认不进入公开消费链路
- agent/search/RAG 后续必须同时遵守 `status + visibility + ai_access`

### 2.4 服务内分层

当前正式实现层次为：

```text
server -> logic -> service -> repo -> entity
```

职责边界：

- `server`：gRPC server 方法分发
- `logic`：RPC transport 适配、metadata 提取、错误映射
- `service`：用例编排、事务边界、业务校验、权限裁决
- `repo`：GORM 持久化访问和行锁封装
- `entity`：GORM 表结构映射

## 3. 当前已落地能力

当前 content 已落地：

- `CreateContent`
- `UpdateContent`
- `GetContent`
- `ListStudioContents`
- `ArchiveContent`
- `ListContentRevisions`
- `GetContentRevision`
- `CreateContentRelation`
- `DeleteContentRelation`
- `ListContentRelations`
- `CreateTag`
- `UpdateTag`
- `DeleteTag`
- `ListTags`
- `ListPublicContents`
- `GetPublicContentBySlug`
- `Ping`

当前数据库已落地：

- `content.items`
- `content.revisions`
- `content.relations`
- `content.tags`
- `content.content_tags`

当前安全与一致性规则：

- Studio service 层要求 `role=admin`
- 新内容默认 `draft + private + ai_access=denied`
- `body_json` 非空时必须是合法 JSON
- `metadata_json` 非空时必须是合法 JSON
- 内容关系只允许出边管理，不允许自关联
- 删除已绑定 tag 返回 `CodeContentTagInUse`
- `content_tags.tag_id` 使用 `ON DELETE RESTRICT`

## 4. 后续能力顺序

### 4.1 content events

下一优先级实现内容事件。

目标：

- 在内容创建、更新、归档、状态变化、可见性变化、AI 访问变化时发布事件
- 为 `indexer`、`search`、`realtime` 提供异步输入
- 通过 `pkg/mq` 抽象 RabbitMQ，业务层只依赖 publisher 接口

### 4.2 search / indexer

事件稳定后实现 search/indexer。

目标：

- `indexer` 消费 content 事件
- 生成 `search_document`
- 公开搜索只索引 `published + public`
- 后续再扩展 member search 和 agent search

### 4.3 attachments / comments

附件和评论在 content 主体、关系和事件稳定后补齐。

目标：

- 附件支持内容封面、正文资源、源材料
- 评论支持公开内容互动
- 评论管理仍由 `admin` 裁决，发表评论后续允许 `member`

## 5. 当前不做

第一阶段暂不做：

- 多 owner / 多协作者编辑权限
- 内容级 ACL
- 字段级权限
- review TBAC
- agent client 细粒度授权
- gateway 聚合式业务编排
