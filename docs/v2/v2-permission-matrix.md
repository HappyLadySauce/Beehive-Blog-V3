# Beehive Blog v2 权限矩阵设计

## 1. 目标

本文件用于把 v2 第一阶段的权限模型从“概念”落成“可执行规则”。

重点解决：

- 谁能看什么
- 谁能改什么
- AI 能读什么
- 搜索能索引什么
- 哪些内容能进入公开站

## 2. 权限模型回顾

第一阶段统一采用 4 个维度：

- `role`
- `status`
- `visibility`
- `ai_access`

## 3. 角色定义

### 3.1 guest

未登录访客。

### 3.2 member

已注册并登录的普通用户。

### 3.3 owner

平台拥有者，拥有全部后台和审阅权限。

### 3.4 agent

系统授权的智能体主体。

说明：

- 第一阶段 `agent` 在权限上视为统一主体
- 具体调用来源通过 `agent_client` 审计

## 4. 内容状态定义

- `draft`
- `review`
- `published`
- `archived`

## 5. 可见性定义

- `public`
- `member`
- `private`

## 6. AI 访问定义

- `allowed`
- `denied`

## 7. 基础判定原则

权限判定建议按下面顺序执行：

1. 先判断主体身份是否存在
2. 再判断内容状态
3. 再判断内容可见性
4. 最后判断 `ai_access`

这意味着：

- `visibility` 决定“人能不能看”
- `ai_access` 决定“AI 能不能读”
- `status` 决定“内容当前是否进入消费链路”

## 8. 公开消费权限矩阵

## 8.1 guest 访问内容

| status | visibility | 是否可见 |
|------|------|------|
| draft | public/member/private | 否 |
| review | public/member/private | 否 |
| published | public | 是 |
| published | member | 否 |
| published | private | 否 |
| archived | public/member/private | 否，默认不开放 |

说明：

- 第一阶段不建议让 `archived` 内容默认继续公开展示
- 如后续需要，可增加归档页策略

## 8.2 member 访问内容

| status | visibility | 是否可见 |
|------|------|------|
| draft | public/member/private | 否 |
| review | public/member/private | 否 |
| published | public | 是 |
| published | member | 是 |
| published | private | 否 |
| archived | public/member/private | 否，默认不开放 |

## 8.3 owner 访问内容

| status | visibility | 是否可见 |
|------|------|------|
| draft | public/member/private | 是 |
| review | public/member/private | 是 |
| published | public/member/private | 是 |
| archived | public/member/private | 是 |

## 8.4 agent 读取内容

Agent 的读取不直接继承人类可见性，需要同时满足：

- 内容状态允许被读取
- 内容可见性允许进入 AI 范围
- `ai_access=allowed`

建议第一阶段规则如下：

| status | visibility | ai_access | agent 是否可读 |
|------|------|------|------|
| draft | public/member/private | allowed/denied | 否 |
| review | public/member/private | allowed/denied | 否 |
| published | public | allowed | 是 |
| published | public | denied | 否 |
| published | member | allowed | 是 |
| published | member | denied | 否 |
| published | private | allowed | 默认否 |
| published | private | denied | 否 |
| archived | any | any | 第一阶段默认否 |

### 结论

第一阶段建议：

- `private` 内容即使 `ai_access=allowed`，也默认不给外部 agent 读取
- 如未来确实需要“仅 AI 可读”的私密知识，可以在第二阶段扩展

这一步是有意收紧，避免把私密内容泄露给外部智能体。

## 9. 搜索索引权限矩阵

## 9.1 公开搜索索引

只有满足以下条件的内容才进入公开搜索：

- `status=published`
- `visibility=public`

并且：

- 私密经历绝不能进入公开搜索
- `draft/review/private/member` 内容不得进入公开搜索

## 9.2 会员搜索索引

登录用户搜索可命中：

- `published + public`
- `published + member`

不应命中：

- `draft`
- `review`
- `private`

## 9.3 owner 搜索

owner 在 Studio 搜索中可以看到：

- 全部内容

但建议搜索结果显式标记：

- status
- visibility
- ai_access

## 9.4 agent 检索

Agent 检索只能返回：

- `published + public + ai_access=allowed`
- `published + member + ai_access=allowed`

不返回：

- `draft`
- `review`
- `private`
- `ai_access=denied`

## 10. 操作权限矩阵

## 10.1 内容创建

| 操作 | guest | member | owner | agent |
|------|------|------|------|------|
| 创建文章 | 否 | 否 | 是 | 否 |
| 创建笔记 | 否 | 否 | 是 | 否 |
| 创建项目 | 否 | 否 | 是 | 否 |
| 创建经历 | 否 | 否 | 是 | 否 |
| 创建时间线事件 | 否 | 否 | 是 | 否 |
| 创建页面 | 否 | 否 | 是 | 否 |

说明：

- 第一阶段内容生产者只有 owner
- agent 不直接创建正式内容，只能生成草稿输出

## 10.2 内容编辑

| 操作 | guest | member | owner | agent |
|------|------|------|------|------|
| 编辑内容 | 否 | 否 | 是 | 否 |
| 删除内容 | 否 | 否 | 是 | 否 |
| 修改状态 | 否 | 否 | 是 | 否 |
| 修改可见性 | 否 | 否 | 是 | 否 |
| 修改 AI 访问 | 否 | 否 | 是 | 否 |

## 10.3 AI 输出相关

| 操作 | guest | member | owner | agent |
|------|------|------|------|------|
| 发起 AI 摘要 | 否 | 否 | 是 | 否 |
| 发起 AI 草稿生成 | 否 | 否 | 是 | 否 |
| 提交 AI 输出进入审阅 | 否 | 否 | 是 | 是 |
| 审核 AI 输出 | 否 | 否 | 是 | 否 |
| 发布 AI 输出 | 否 | 否 | 是 | 否 |

说明：

- agent 可以生成输出并提交审阅
- 最终审核和发布只归 owner

## 10.4 评论权限

| 操作 | guest | member | owner | agent |
|------|------|------|------|------|
| 查看公开评论 | 是 | 是 | 是 | 否 |
| 发表评论 | 否 | 是 | 是 | 否 |
| 删除自己评论 | 否 | 第二阶段可支持 | 是 | 否 |
| 隐藏评论 | 否 | 否 | 是 | 否 |
| 删除评论 | 否 | 否 | 是 | 否 |

第一阶段评论策略：

- member 可直接发表评论
- 默认不做人审阻塞
- 可挂接简单审核接口或风控检查

## 10.5 用户权限

| 操作 | guest | member | owner | agent |
|------|------|------|------|------|
| 注册 | 是 | 否 | 否 | 否 |
| 登录 | 是 | 是 | 是 | 否 |
| 查看自己的资料 | 否 | 是 | 是 | 否 |
| 修改自己的资料 | 否 | 是 | 是 | 否 |
| 管理用户 | 否 | 否 | 是 | 否 |

## 11. 默认值策略

为避免权限混乱，建议所有新内容创建时使用统一默认值。

## 11.1 owner 手工创建内容默认值

- `status=draft`
- `visibility=private`
- `ai_access=denied`

## 11.2 AI 生成内容默认值

- `status=draft`
- `visibility=private`
- `ai_access=denied`

## 11.3 发布后的默认值建议

### article

- 默认建议 `published + public + ai_access=allowed`

### note

- 默认建议 `draft/private`
  或 `published/member/allowed`

### experience

- 默认建议 `draft/private/denied`

### timeline_event

- 默认建议跟随所属经历，不自动公开

## 12. 权限检查落点建议

## 12.1 gateway

负责：

- 登录态校验
- token 基础解析

不负责：

- 内容级权限决策

## 12.2 content-service

负责：

- 内容是否可读
- 内容是否可编辑
- 评论是否允许挂载

## 12.3 search-service

负责：

- 索引时剔除不该进入搜索的内容
- 查询时二次过滤

## 12.4 agent-service

负责：

- AI 读取前再校验 `status + visibility + ai_access`
- 防止越权读取 member/private 内容

## 13. 第一阶段故意不支持的权限能力

为了控制复杂度，第一阶段不建议做：

- 多 owner / 多协作者编辑权限
- 细粒度字段级权限
- 多空间隔离
- 按 agent client 单独授权不同内容范围
- “仅 AI 可读但 owner 不公开”的复杂模式

## 14. 推荐实现顺序

1. 先在 `content_items` 落 `status + visibility + ai_access`
2. 先实现内容读取权限判定函数
3. 再实现搜索索引过滤策略
4. 最后实现 agent 读取权限判定

## 15. 当前结论

v2 第一阶段的权限模型可以收束为：

**人类可见性由 `visibility` 控制，内容生命周期由 `status` 控制，AI 读取由 `ai_access` 控制，最终由各服务按统一矩阵执行。**
