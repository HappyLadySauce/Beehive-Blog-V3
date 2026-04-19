# Beehive Blog v2.0 架构设计草案

## 1. 背景

Beehive Blog v1 当前已经具备以下价值：

- 提供了可运行的 Go 后端、Hexo 前台、React 管理后台
- 积累了文章、标签、分类、评论、附件、鉴权等业务实现经验
- 形成了初步的 Hexo 同步设计和后台内容管理能力

但从 v2 的目标来看，v1 的整体结构仍然有明显限制：

- 内容模型以“博客文章”为中心，不足以承载“个人知识库 + 经历地图 + AI 协作”
- Hexo 参与了过多内容层职责，导致展示层和内容层边界不够清晰
- AI 应用接入尚未成为一等公民，缺少稳定的读写协议与审核流
- 前端能力主要服务博客展示，不适合构建知识检索、关系浏览、时间线、AI 审阅等复杂体验

因此，v2 不再以“优化 Hexo 博客”作为目标，而是以“构建个人知识中台”为核心方向。

## 2. v2 愿景

Beehive Blog v2.0 的目标不是单纯的博客系统，而是一个同时服务于人类与智能体的个人知识平台：

- 对个人：沉淀经历、认知、项目、文章、作品与人生轨迹
- 对访客：可检索、可学习、可引用、可参考
- 对智能体：可标准化读取、理解、归纳与辅助生成
- 对创作流程：支持人类与 AI 的共同写作、共同整理、共同演进

一句话定义：

**Beehive Blog v2.0 = 个人知识中台 + 多终端前端 + AI 双向协作系统**

## 3. 核心设计原则

### 3.1 后端是唯一可信内容源

所有内容真相都应归属于后端的数据模型与内容服务。

- Markdown 不是业务真相，只是一种导入导出格式
- Hexo 不是内容中心，只是一个可选展示适配器
- 未来新增前端、MCP、Skill、Hermes、OpenClaw 时，不应改动核心内容模型

### 3.2 AI 接入必须是一等公民

AI 读写能力不是外挂功能，而是架构级能力。

- 需要有独立的 Agent API / MCP Server / Skill Output 流程
- AI 写入必须具备来源、权限、审核、回滚、版本记录
- AI 读取必须拿到结构化、可信、可追溯的数据

### 3.3 内容模型必须从“文章”升级为“知识实体”

v2 的核心对象不应只有 Article，还要支持：

- Experience：个人经历
- Event：时间线事件
- Project：项目
- Note：知识笔记
- Insight：反思与认知
- Artifact：附件、截图、代码片段、作品
- Person：人物
- Relation：实体关系

### 3.4 展示层可替换

公开博客、个人知识后台、移动端、AI 对话端都只是不同终端。

- 官方 Web 前端可以自研
- 公开站可以 SSR / SSG
- Hexo 可保留为历史兼容或静态导出目标
- 任一展示层都不应绑定核心业务规则

### 3.5 先可审阅，再可自动化

AI 产出不能直接默认发布。

- 先生成草稿
- 再进行人工审阅或规则校验
- 最后发布到公开知识空间

## 4. v2 总体架构

```text
                +---------------------------+
                |    Human Creator/Admin    |
                +-------------+-------------+
                              |
                              v
                +---------------------------+
                |     Web Console (v2)      |
                |  内容编辑 / 审阅 / 运营   |
                +-------------+-------------+
                              |
                              v
+----------------+   +---------------------------+   +-------------------+
|  Agent Client  |-->|      API Gateway          |<--|   Public Web App   |
| Hermes/OpenClaw|   | REST / Auth / RateLimit   |   | SSR/SSG/Hybrid UI  |
+----------------+   +-------------+-------------+   +-------------------+
                              |
                +-------------+-------------+
                |     Domain Services        |
                | Content / Knowledge / AI   |
                | Publish / Search / Review  |
                +------+------+------+-------+
                       |      |      |
                       v      v      v
                 +--------+ +-------+ +------------------+
                 |  PGSQL | | Redis | | Search / Vector  |
                 +--------+ +-------+ +------------------+
                       |
                       v
             +------------------------+
             | Adapter Layer          |
             | Hexo / Markdown / RSS  |
             | MCP / Skill / Export   |
             +------------------------+
```

## 5. 推荐模块拆分

### 5.1 Core API

职责：

- 提供认证、权限、内容查询、编辑、发布、审阅、搜索等统一 API
- 为 Web 前端、管理后台、Agent 客户端提供统一入口

建议拆分：

- `identity-service`
- `content-service`
- `knowledge-service`
- `publish-service`
- `review-service`
- `search-service`
- `agent-service`

如果前期不希望过早微服务化，可以先保留单体仓库，但在代码层明确模块边界。

### 5.2 Content Service

职责：

- 管理文章、页面、笔记、经历、项目、事件等内容实体
- 负责草稿、版本、状态流转、发布时间、可见性等规则

建议所有内容统一抽象为：

- `content_item`
- `content_revision`
- `content_block`
- `content_relation`
- `content_metadata`

在抽象层之上再做具体类型扩展。

### 5.3 Knowledge Service

职责：

- 从内容实体中生成知识视图
- 支持时间线、关系图、主题聚合、经历链路、个人画像
- 为智能体提供结构化知识检索接口

### 5.4 Agent Service

职责：

- 提供面向智能体的标准接入层
- 暴露 MCP 能力
- 管理 Skill 输出写入
- 控制 AI 草稿生成、归档、审阅、发布流

建议支持两类接口：

- Read API：检索知识、读取上下文、获取可引用材料
- Write API：提交草稿、生成摘要、更新标签建议、创建待审内容

### 5.5 Publish Service

职责：

- 将已发布内容输出到不同目标
- 支持站内 Web 渲染
- 支持静态导出
- 支持 Hexo 兼容导出
- 支持 RSS / Sitemap / JSON Feed / Markdown Export

这层的关键是把“内容中心”和“输出目标”彻底解耦。

### 5.6 Search Service

职责：

- 关键词搜索
- 标签、分类、主题筛选
- 时间范围检索
- 相似内容推荐
- 未来接入全文索引和向量检索

建议分两阶段：

- 第一阶段：PostgreSQL 全文检索 + 业务过滤
- 第二阶段：外接 Meilisearch / Elasticsearch / OpenSearch + Vector Store

## 6. 新的数据建模方向

### 6.1 内容主实体

建议引入统一内容表，而不是把所有能力都硬塞到 `articles`。

核心字段建议：

- `id`
- `type`
- `title`
- `slug`
- `status`
- `visibility`
- `summary`
- `author_id`
- `author_type`
- `source_type`
- `source_ref`
- `published_at`
- `created_at`
- `updated_at`
- `current_revision_id`

其中：

- `type`: article / note / page / experience / project / event / insight
- `status`: draft / review / scheduled / published / archived / deleted
- `visibility`: public / unlisted / private / agent_only
- `author_type`: human / agent / hybrid
- `source_type`: manual / import_markdown / import_hexo / mcp / skill / agent

### 6.2 版本模型

需要显式版本表，不能只靠当前文章内容覆盖。

- `content_revision`
- `revision_number`
- `editor_id`
- `editor_type`
- `change_summary`
- `body_markdown`
- `body_json`
- `citations_json`
- `created_at`

这样才能支持：

- AI 草稿与人工改稿并存
- 差异对比
- 回滚
- 审计

### 6.3 关系模型

建议新增通用关系表：

- `from_type`
- `from_id`
- `relation_type`
- `to_type`
- `to_id`
- `weight`
- `metadata`

关系例子：

- 文章属于某段经历
- 项目使用某项技术
- 某条认知来自某个事件
- 一篇文章引用另一篇文章
- 某个附件属于某个项目

### 6.4 AI 草稿与审阅

建议单独设计：

- `agent_task`
- `agent_output`
- `review_task`
- `review_decision`

最少要能记录：

- 哪个智能体生成的
- 使用了什么上下文
- 生成目标是什么
- 输出到哪个内容实体
- 当前审核状态是什么
- 谁批准发布

## 7. 面向 AI 的标准能力设计

### 7.1 AI 读能力

智能体需要读取的不是简单 HTML，而是结构化知识。

建议提供：

- `searchKnowledge`
- `getContentById`
- `getContentBySlug`
- `getTimeline`
- `getExperienceGraph`
- `getRelatedContent`
- `getRevisionHistory`
- `getCitationContext`

输出应包含：

- 标题
- 摘要
- 正文
- 标签
- 类型
- 时间
- 引用来源
- 关联实体
- 可信度和可见性

### 7.2 AI 写能力

智能体写入只进入草稿区，不直接发布。

建议提供：

- `createDraft`
- `appendDraftSection`
- `suggestTags`
- `linkRelatedEntities`
- `submitForReview`
- `createExperienceSummary`
- `generateWeeklyDigest`

### 7.3 MCP 支持

v2 可直接提供本地或远程 MCP Server。

建议暴露工具：

- `search_content`
- `read_content`
- `write_draft`
- `list_taxonomy`
- `list_projects`
- `list_experiences`
- `submit_review_request`

资源建议：

- `content://items/{id}`
- `content://timeline`
- `content://projects/{slug}`
- `content://experiences/{slug}`
- `content://tags`

## 8. 前端重构建议

## 8.1 结论

建议 v2 跳出 Hexo 作为主前端框架，自研新的 Web 前端。

原因：

- 你们要做的是“知识平台”，不是“传统静态博客”
- 需要复杂状态、权限、检索、关系视图、AI 审阅流
- 这些场景更适合现代 React 前端或 SSR/Hybrid 框架

Hexo 在 v2 中建议降级为：

- 兼容导出目标
- 历史内容迁移来源
- 可选静态镜像发布器

### 8.2 前端建议分层

建议拆成两个 Web 应用：

- `web/public`
  - 对外公开站点
  - 文章阅读、时间线、项目页、知识地图、搜索
- `web/studio`
  - 内容编辑台
  - AI 草稿审阅、版本管理、知识整理、发布管理

### 8.3 前端技术建议

如果以长期演进为目标，建议优先考虑：

- `Next.js` 或 `Remix` 作为公开站
- React + TypeScript
- SSR + SSG + 动态数据混合
- 统一消费 Go API

选择理由：

- SEO、内容站点能力成熟
- 动态交互比 Hexo 强得多
- 对搜索、鉴权、个性化页面更友好

如果团队明确想保持 Go 为主，也可以考虑：

- Go 提供 API
- 前端单独 React/SSR 工程
- 最终通过 nginx/caddy 聚合部署

## 9. 建议的仓库形态

建议 v2 采用 monorepo，但明确边界。

```text
beehive-v2/
  apps/
    api/                # Go API
    mcp-server/         # 面向智能体的 MCP 服务
    web-public/         # 对外公开站
    web-studio/         # 管理后台 / 编辑工作台
    worker/             # 异步任务：索引、导出、摘要、调度
  packages/
    domain/             # 共享领域模型
    sdk/                # TS/Go SDK
    prompt-templates/   # AI 模板
    adapters/           # hexo / markdown / rss / export
  docs/
    architecture/
    roadmap/
    api/
    adr/
```

## 10. 分阶段建设路线

### Phase 0：冻结 v1，作为参考基线

- 停止在 v1 上继续扩展核心能力
- 仅保留 bugfix 或数据导出支持
- 盘点可复用模块：鉴权、附件、文章版本、标签分类、同步经验

### Phase 1：搭建 v2 内容中台

目标：

- 定义新数据模型
- 初始化 API 工程
- 完成认证、内容实体、版本、草稿、审阅基础能力

交付：

- 统一内容表
- 版本表
- 关系表
- 审阅流
- 基础管理 API

### Phase 2：接入新的公开前端与工作台

目标：

- 替换 Hexo 的主展示职责
- 建立新的公开站和编辑工作台

交付：

- 公开内容页
- 项目页
- 时间线页
- 知识检索页
- 草稿审阅页

### Phase 3：补齐 AI 读写标准接口

目标：

- 构建 MCP Server
- 支持 Skill 写入草稿
- 支持基于知识库的智能体学习与引用

交付：

- MCP 工具
- Agent Output 草稿流
- 引用追踪
- AI 写作审阅流程

### Phase 4：多输出通道

目标：

- 支持静态导出与外部发布

交付：

- Hexo Adapter
- Markdown Export
- RSS / Sitemap / JSON Feed
- 第三方知识平台同步能力

## 11. 与 v1 的关系

v1 不再作为继续演化的基础框架，而是作为以下三类参考：

- 业务规则参考
- 数据迁移来源
- 可复用模块来源

优先复用的内容建议：

- 鉴权模型与部分中间件思想
- 附件存储与上传处理经验
- 数据库命名和内容管理经验
- Hexo 同步经验中的适配思路

不建议直接沿用的内容：

- 以 Hexo 为中心的内容组织方式
- 以文章为唯一主实体的业务边界
- 前台和知识系统耦合的渲染模式

## 12. v2 成功标准

如果 v2 设计正确，最终应具备以下特征：

- 人类可以把博客、经历、项目、反思统一沉淀进一个知识空间
- AI 可以按标准协议读取知识，并输出草稿而不是污染正式内容
- 公开站只是内容的一种表现形式，不再限制业务设计
- 内容具备结构化、可检索、可追踪、可演化能力
- 系统能支持“个人成长记录”与“AI 协作创作”长期共存

## 13. 当前建议的下一步

建议接下来紧接着完成三份文档：

1. `v2-domain-model.md`
   - 设计实体、字段、关系、状态机
2. `v2-api-design.md`
   - 设计内容、审阅、搜索、Agent、MCP 接口
3. `v2-migration-plan.md`
   - 设计从 v1 数据迁移到 v2 的策略

在这三份文档完成之前，不建议直接进入编码阶段。
