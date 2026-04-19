# Beehive Blog v2.0 搜索、索引与 RAG 设计草案

## 1. 这份文档解决什么问题

v2 已经明确要走“个人知识中台 + AI 协作平台”的方向。

这意味着系统不能只支持：

- 浏览文章
- 按标签筛选
- 简单全文搜索

而要进一步支持：

- 把个人博客、项目、经历、笔记组织成可检索知识库
- 让人类可以快速搜索、浏览、理解自己的内容
- 让智能体可以按标准协议读取知识，而不是盲目爬页面
- 让 AI 基于检索结果生成更短、更准、更可追溯的回答或草稿

因此，v2 需要明确设计三层能力：

1. 索引
2. 检索
3. RAG

## 2. 先说结论：我们要不要做 RAG

### 2.1 简短结论

**要做，但不建议一上来就做重型 RAG。**

更准确地说：

- 第一阶段必须做“索引 + 搜索”
- 第二阶段做“混合检索”
- 第三阶段再做“RAG 问答与 AI 写作辅助”

### 2.2 为什么不是先上 RAG

RAG 的效果高度依赖前面的基础设施：

- 内容是否结构化
- 索引是否完整
- 检索是否准确
- 引用来源是否可靠
- 权限边界是否明确

如果这些没做好，RAG 只会把错误内容更流畅地说出来。

### 2.3 RAG 到底是什么

RAG 可以简单理解为：

1. 先从你的知识库里检索相关内容
2. 再把检索结果交给大模型生成回答、摘要或草稿

所以：

- RAG 不是训练模型
- RAG 不是内容压缩工具
- RAG 也不是搜索引擎替代品

它更像是：

**“搜索 + 上下文拼装 + 生成”**

### 2.4 关于“压缩得更加短小精悍”

这里需要纠正一个关键认知：

**RAG 本身不负责压缩原始内容。**

RAG 能带来的改进是：

- 在回答问题时，只取和当前问题相关的少量片段
- 然后由模型把这些片段整理成更短的结果
- 同时附带引用来源，避免“瞎编”

所以真正的链路是：

**原始内容 -> 切片与索引 -> 检索命中 -> 模型总结 -> 输出简洁答案**

不是：

**原始内容 -> RAG -> 永久压缩**

如果你们想要“长期更短更精炼”的内容沉淀，还需要单独设计：

- 摘要生成
- 重点提炼
- 周报/阶段总结
- 文章与经历的多层摘要

这属于“知识提炼功能”，不是 RAG 的全部。

## 3. v2 中搜索、索引、RAG 的关系

```text
原始内容录入
  -> 内容标准化
  -> 内容切片
  -> 元数据提取
  -> 关键词索引
  -> 向量索引
  -> 混合检索
  -> RAG 上下文组装
  -> AI 回答 / AI 草稿 / 人类搜索结果页
```

更具体一点：

```text
Content Item
  -> Revision
  -> Chunk
  -> Search Document
  -> Embedding
  -> Retrieval Result
  -> Citation Bundle
  -> LLM Output
```

## 4. 我们建议的总体方案

### 4.1 大方向

建议采用：

- 结构化内容存储：PostgreSQL
- 基础关键词检索：PostgreSQL Full Text Search
- 中大型索引引擎：Meilisearch 或 Elasticsearch
- 向量检索：pgvector
- 混合检索：关键词分数 + 向量分数融合
- RAG：在后端或独立 worker 中实现
- 外部智能体接入：MCP Server + Agent API

### 4.2 为什么这样选

对你们当前阶段最合适的原因有三点：

1. 架构复杂度可控
2. 和 Go 后端、PostgreSQL 技术栈兼容
3. 方便先把知识库能力做起来，而不是过早引入过重基础设施

不过，如果你们从一开始就明确要做以下场景，也可以直接把独立搜索引擎纳入一期设计：

- 更强的搜索体验
- 更快的模糊搜索与联想词
- 更复杂的排序调优
- 面向公开站和 Studio 的统一搜索中心
- 为后续 AI 检索提供更稳定的召回基础

## 4.3 搜索引擎选型建议

### 方案 A：PostgreSQL Full Text Search

优点：

- 技术栈最简单
- 部署和运维成本低
- 适合早期快速落地
- 与业务过滤条件结合自然

不足：

- 搜索体验上限有限
- 联想词、拼写容错、排序调优能力一般
- 不适合作为长期复杂搜索中心

适用阶段：

- v2 早期原型
- 内容量较小
- 团队希望优先验证内容模型

### 方案 B：Meilisearch

优点：

- 接入简单，开发体验好
- 默认搜索体验通常比数据库全文检索更好
- 联想、容错、排序、过滤能力比较实用
- 很适合博客、文档、知识库类站点

不足：

- 复杂度高于直接用 PostgreSQL
- 向量与混合检索能力可作为后续能力考虑，但整体生态和灵活度不如 Elasticsearch 体系
- 对非常复杂的企业级检索场景扩展性相对有限

适用阶段：

- 你们希望尽快获得“像样的搜索体验”
- 团队希望先把搜索做得好用，而不是先研究复杂搜索平台
- 内容主要是博客、知识卡片、项目记录、时间线

### 方案 C：Elasticsearch

优点：

- 检索能力强，生态成熟
- 适合复杂查询、复杂排序、多索引策略
- 更适合中长期承担搜索中心角色
- 对大规模内容、复杂聚合、可观测性和调优更友好

不足：

- 部署和运维复杂度更高
- 建设成本高于 Meilisearch
- 对当前阶段可能偏重

适用阶段：

- 你们明确要把搜索做成平台级能力
- 会有复杂检索、聚合分析、推荐和 AI 检索需求
- 团队接受更高的基础设施复杂度

### 当前推荐

结合你们的目标，我建议文档里的推荐调整为：

- 近期优先方案：`Meilisearch + PostgreSQL`
- 中长期增强方案：`Elasticsearch + Vector Store`

也就是说：

- 如果你们希望尽快做出体验不错的知识库搜索，用 Meilisearch 更合适
- 如果你们已经明确未来会做复杂检索平台，可以直接以 Elasticsearch 为目标设计接口边界

不管选哪种引擎，内容主数据仍应留在 PostgreSQL，搜索引擎只承担索引与检索职责。

## 5. 系统功能拆解

## 5.1 内容标准化模块

职责：

- 接收文章、笔记、项目、经历、事件等内容
- 清洗正文、标题、标签、时间、关联关系
- 生成统一的内部内容结构

输入：

- 人工编辑内容
- Markdown 导入
- v1 迁移数据
- MCP / Skill 写入草稿

输出：

- `content_item`
- `content_revision`
- `content_relation`

这是全部检索和 RAG 的前提。

## 5.2 内容切片模块

职责：

- 把长文拆成更适合检索的片段
- 保留每个片段的上下文、顺序和来源

建议切片粒度：

- 标题
- 一级/二级小节
- 约 300-800 中文字的正文块
- 代码块、引用块、时间线块单独切分

切片必须保留以下元数据：

- `content_id`
- `revision_id`
- `chunk_id`
- `chunk_type`
- `position`
- `heading_path`
- `visibility`
- `author_type`
- `published_at`
- `tags`
- `project_id`
- `experience_id`

## 5.3 关键词索引模块

职责：

- 提供精确关键词搜索
- 支持标题权重高于正文
- 支持标签、时间、类型、可见性过滤

建议实现：

- PostgreSQL `tsvector`
- GIN 索引
- 标题/摘要/正文分权重
- 或同步构建 Meilisearch / Elasticsearch 文档索引

这个能力最适合解决：

- 你记得关键字，但不记得文章名
- 需要精确搜人名、技术名、项目名、事件名
- 需要做后台搜索页和公开站搜索页

## 5.4 向量索引模块

职责：

- 支持语义检索
- 支持“意思相近但关键词不同”的召回
- 为 RAG 提供候选上下文

建议实现：

- 为 `chunk` 生成 embedding
- 存储到 PostgreSQL `pgvector`
- 初期使用 HNSW

适合解决：

- “我以前写过类似内容，但我不记得原词”
- “找和这个经历相关的反思”
- “找和这个项目思路接近的文章”

## 5.5 混合检索模块

职责：

- 将关键词检索和向量检索结果融合
- 兼顾精确命中与语义召回

建议排序因子：

- 关键词分数
- 向量相似度分数
- 内容类型权重
- 发布时间衰减
- 手动权重
- 质量分或可信度分

建议实现形态：

- PostgreSQL FTS + pgvector 融合
- Meilisearch 关键词召回 + pgvector 二次融合
- Elasticsearch 检索召回 + 向量召回 + 重排

适用场景：

- 面向用户的全站搜索
- 面向智能体的知识检索
- 面向写作助手的“找相关材料”

## 5.6 摘要与知识提炼模块

职责：

- 生成短摘要、长摘要、要点列表
- 为搜索结果页和智能体调用提供更短上下文

建议产物：

- `summary_short`
- `summary_medium`
- `summary_long`
- `key_points`
- `entities`
- `citations`

注意：

这部分和 RAG 有关系，但不是同一件事。

它更像是“预处理压缩层”。

## 5.7 RAG 问答模块

职责：

- 接收用户问题
- 检索相关片段
- 组装上下文
- 调用模型生成回答
- 返回引用来源

推荐只服务以下场景：

- 站内“问我的知识库”
- 管理后台“根据历史内容辅助写作”
- MCP 给外部智能体提供知识问答能力

不建议它直接承担：

- 正式内容发布
- 全站搜索主入口

## 5.8 AI 草稿生成模块

职责：

- 基于检索结果生成文章草稿、经历总结、阶段复盘
- 生成后进入待审状态

输入：

- 用户命题
- 指定内容范围
- 检索命中的上下文

输出：

- draft revision
- 引用来源
- 变更说明

## 5.9 引用追踪模块

职责：

- 记录 AI 回答或草稿引用了哪些 chunk
- 给用户可追溯证据
- 方便后续校对和回滚

这一步非常关键。

没有引用追踪，RAG 的可信度会明显下降。

## 5.10 MCP / Agent 接入模块

职责：

- 通过 MCP 暴露搜索、读取、草稿写入等能力
- 让 Hermes、OpenClaw 等智能体应用标准化接入

推荐能力：

- `search_content`
- `read_content`
- `read_chunk`
- `ask_knowledge`
- `write_draft`
- `submit_review`

## 6. 功能之间的联系关系

## 6.1 关系总览

```text
内容录入
  -> 内容标准化
  -> 内容切片
  -> 关键词索引
  -> 向量索引
  -> 混合检索
  -> 摘要提炼
  -> RAG
  -> AI 草稿
  -> 审阅
  -> 发布
```

## 6.2 依赖关系

### A. 没有内容标准化，就没有高质量索引

因为：

- 标题、摘要、标签、实体信息缺失
- 检索只能扫大段正文
- AI 也拿不到结构化上下文

### B. 没有内容切片，就没有可用的 RAG

因为：

- 大模型不能每次吃整篇文章
- 太长的上下文会贵、慢、还不稳定
- 返回引用也难以精确定位

### C. 没有关键词索引，基础搜索体验会很差

因为语义检索虽然强，但对：

- 专有名词
- 错别字
- 精确词汇
- 标签和 ID 类信息

往往不如关键词检索可靠。

### D. 没有向量索引，AI 问答会“只会找字面匹配”

这会导致：

- 相似内容召回不足
- 表达方式变了就搜不到
- 个人经历类内容很难串起来

### E. 没有引用追踪，RAG 结果不可审计

这会直接影响：

- 可信度
- 审阅效率
- 对外展示安全性

## 7. 你们真正需要的不是“是否启用 RAG”，而是分层建设

建议把能力拆成 4 层。

### Layer 1：搜索基础层

目标：

- 先把“能搜到”做好

包含功能：

- 内容标准化
- 关键词索引
- 搜索 API
- 高亮摘要
- 标签/类型/时间过滤

对外能力：

- 网站搜索
- 管理台搜索
- MCP 基础读取

### Layer 2：知识检索层

目标：

- 先把“能搜得准”做好

包含功能：

- 内容切片
- embedding 生成
- 向量索引
- 混合检索
- related content

对外能力：

- 相似内容推荐
- 经历关联内容
- 项目相关资料聚合

### Layer 3：知识提炼层

目标：

- 先把“内容能变短、变精华”做好

包含功能：

- 自动摘要
- 要点提取
- 实体抽取
- 周报/阶段总结

对外能力：

- 搜索结果摘要
- 个人成长周报
- 文章简介

### Layer 4：RAG 协作层

目标：

- 最后再把“AI 真的用起来”做好

包含功能：

- ask_knowledge
- grounded answer
- draft generation
- citation bundle
- review workflow

对外能力：

- 知识问答
- 辅助写作
- MCP 智能体协作

## 8. 推荐的阶段计划

## Phase 1：基础搜索上线

范围：

- PostgreSQL FTS 或 Meilisearch
- 搜索页
- 搜索 API
- 搜索高亮
- 标签、类型、时间过滤

暂时不做：

- 向量检索
- RAG

目的：

- 尽快拿到可用搜索体验
- 验证内容模型和元数据设计是否合理

## Phase 2：混合检索上线

范围：

- chunk 表
- embedding 生成 worker
- pgvector
- hybrid rank
- related content API
- 如果一期已采用 PostgreSQL FTS，可在这一阶段补 Meilisearch 或 Elasticsearch
- 如果一期已采用 Meilisearch，可在这一阶段补向量索引与结果融合

目的：

- 提升检索召回率
- 为后续 RAG 做准备

## Phase 3：知识提炼上线

范围：

- 摘要
- 要点
- 实体抽取
- 个人经历总结

目的：

- 真正实现“内容更短、更精华”
- 让搜索结果和知识卡片更有用

## Phase 4：RAG 问答与 AI 草稿

范围：

- ask_knowledge
- AI 引用回答
- AI 草稿写入
- 审阅流

目的：

- 让 AI 成为知识库协作者
- 而不是直接污染正式内容

## 9. 推荐的数据实体

在原有内容实体之外，建议新增：

### 9.1 search_document

作用：

- 面向检索的聚合文档

建议字段：

- `id`
- `content_id`
- `revision_id`
- `title`
- `summary`
- `body_text`
- `search_vector`
- `language`
- `visibility`
- `published_at`

### 9.2 content_chunk

作用：

- 存储切片结果

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

### 9.3 chunk_embedding

作用：

- 存储 chunk 对应向量

建议字段：

- `chunk_id`
- `model_name`
- `embedding`
- `embedding_dim`
- `created_at`

### 9.4 content_summary

作用：

- 存储多层摘要

建议字段：

- `content_id`
- `revision_id`
- `summary_type`
- `text`
- `source`
- `created_at`

### 9.5 retrieval_log

作用：

- 记录一次检索调用命中了哪些内容

### 9.6 rag_answer_log

作用：

- 记录一次 RAG 回答使用了哪些 chunk 和哪些提示词模板

## 10. 推荐 API 能力

## 10.1 面向用户

- `GET /search`
- `GET /search/suggest`
- `GET /content/:id/related`
- `GET /timeline/search`

## 10.2 面向管理后台

- `POST /admin/index/rebuild`
- `POST /admin/index/rebuild-content/:id`
- `GET /admin/search/debug`
- `GET /admin/rag/answers/:id`

## 10.3 面向 AI / MCP

- `search_content`
- `read_content`
- `read_chunks`
- `ask_knowledge`
- `create_draft_from_query`
- `summarize_content`

## 11. 前端需要新增哪些页面

### 11.1 公开站

- 全站搜索页
- 搜索结果详情页
- 相关内容推荐卡片
- 个人经历知识地图页
- “问知识库”页

### 11.2 Studio 工作台

- 索引状态页
- 检索调试页
- AI 回答记录页
- AI 草稿审阅页
- 内容摘要管理页

## 12. 建议的技术取舍

## 12.1 第一阶段推荐

- PostgreSQL Full Text Search
- 或 Meilisearch
- GIN
- `ts_rank`
- `ts_headline`
- pgvector

原因：

- PostgreSQL 方案最轻
- Meilisearch 方案更偏向“搜索体验优先”
- 二者都适合早期建设

## 12.2 暂时不建议一开始就引入

- 单独的重量级搜索集群
- 复杂向量数据库集群
- 很重的 Agent 编排平台

原因：

- 当前最缺的不是基础设施，而是明确的数据模型和能力边界

## 12.3 未来升级方向

当内容量和检索流量上来后，再考虑：

- Elasticsearch
- OpenSearch 混合检索
- 专门的向量服务
- 更复杂的重排模型
- 查询质量评估平台

## 13. 风险与边界

### 13.1 最大风险不是技术，而是数据质量

如果原始内容：

- 标签乱
- 标题弱
- 摘要缺失
- 结构混乱

那么索引和 RAG 效果都会差。

### 13.2 AI 输出必须带引用

否则你们后面很难判断：

- 它是从哪篇文章学来的
- 有没有误解原意
- 是否适合公开发布

### 13.3 权限必须前置

RAG 和检索都必须遵循内容可见性：

- `public`
- `private`
- `draft`
- `agent_only`

不能因为做了索引，就把私密内容泄露给公开查询或外部智能体。

## 14. 最终建议

对于 Beehive Blog v2，我建议采用下面的路线：

1. 先做内容模型与元数据标准化
2. 再做 PostgreSQL 全文索引或 Meilisearch 基础搜索
3. 然后做 chunk、embedding、pgvector 和混合检索
4. 内容规模和搜索复杂度上来后，视情况升级到 Elasticsearch
5. 再做摘要提炼能力
6. 最后做 RAG 问答和 AI 草稿协作

换句话说：

**需要 RAG，但应该把它放在搜索体系的上层，而不是拿它代替搜索体系本身。**

## 15. 当前建议的下一步

建议接下来继续补两份文档：

1. `v2-domain-model.md`
   - 把搜索实体、chunk、summary、rag log 一并纳入正式数据模型
2. `v2-feature-map.md`
   - 把“内容、搜索、AI、发布、审阅”之间的功能关系画成完整矩阵

在这两份文档定下来之前，不建议直接实现 RAG。

## 15.1 当前已实现基线（2026-04-16）

- 已有 `search_documents`、`content_chunks`、`content_summaries`、`rag_answer_logs` 表。
- 已有搜索 RPC：`Query`、`UpsertDocument`、`DeleteDocument`。
- 已有网关能力：公开搜索、Studio 搜索、按内容重建索引、删除索引。
- 内容写操作（创建/更新/状态变更）已异步触发索引更新。

## 16. 参考资料

以下资料用于校验技术判断：

- PostgreSQL Full Text Search 官方文档：
  https://www.postgresql.org/docs/current/textsearch.html
- PostgreSQL 文本搜索索引：
  https://www.postgresql.org/docs/current/textsearch-indexes.html
- PostgreSQL 文本搜索控制、排序与高亮：
  https://www.postgresql.org/docs/current/textsearch-controls.html
- Meilisearch 官方文档：
  https://www.meilisearch.com/docs
- Elasticsearch 官方文档：
  https://www.elastic.co/guide/index.html
- pgvector 官方说明：
  https://pgxn.org/dist/vector/0.8.2/
- OpenSearch Hybrid Search 官方文档：
  https://docs.opensearch.org/2.16/search-plugins/hybrid-search/
- MCP 官方规范：
  https://modelcontextprotocol.io/specification/2025-06-18
- MCP Tools 规范：
  https://modelcontextprotocol.io/specification/2025-03-26/server/tools
