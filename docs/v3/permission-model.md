# Beehive Blog v3 第一阶段权限模型

## 1. 目标

本文件定义 v3 第一阶段的权限模型，作为 `gateway`、`identity`、`content`、后续 `search` 与 `agent` 服务实现授权判断的统一依据。

第一阶段采用 **RBAC + 轻量 ABAC**：

- RBAC：使用 `role` 判断主体基础权限。
- 轻量 ABAC：使用 `status`、`visibility`、`ai_access` 判断内容资源访问范围。
- TBAC：暂不进入第一阶段，仅预留给后续 `review` 审阅工作流。
- OBAC：暂不进入第一阶段，不做逐内容 ACL、协作者或字段级权限。

## 2. 主体角色

v3 第一阶段固定使用以下角色：

- `guest`：未登录访客，不落库。
- `member`：普通登录用户，默认注册角色。
- `admin`：平台管理角色，对应 v2 文档中的 `owner`。

`identity` 负责认证、账号状态和角色真相；业务服务负责资源级授权裁决。

## 3. 内容资源属性

内容权限判断使用以下资源属性：

- `status`：`draft`、`review`、`published`、`archived`
- `visibility`：`public`、`member`、`private`
- `ai_access`：`allowed`、`denied`

判定顺序固定为：

1. 认证主体：是否登录、token 是否有效、账号是否 active。
2. RBAC：角色是否允许进入该业务能力。
3. 内容状态：内容是否处于可消费或可管理状态。
4. 可见性：人类用户是否可读。
5. AI 访问：agent、search、RAG 是否可读。

## 4. 第一阶段规则

### 4.1 Studio 内容管理

- `admin` 可以创建、读取、更新、归档内容。
- `admin` 可以管理标签、读取版本。
- `member` 和 `guest` 不允许访问 Studio 内容管理能力。
- `gateway` 只负责认证和身份透传，最终授权裁决在 `content` 服务。

### 4.2 公开读取

| 主体 | 可读内容 |
| --- | --- |
| `guest` | `published + public` |
| `member` | `published + public`、`published + member` |
| `admin` | 全部内容状态与可见性 |

`archived` 默认不进入公开消费链路。

### 4.3 AI / agent 读取

第一阶段只允许 agent 读取：

- `published + public + ai_access=allowed`
- `published + member + ai_access=allowed`

以下内容默认拒绝外部 agent 读取：

- `draft`
- `review`
- `private`
- `archived`
- 任意 `ai_access=denied` 内容

### 4.4 默认值

新内容默认值：

- `status=draft`
- `visibility=private`
- `ai_access=denied`

发布时由 `admin` 显式决定是否调整为 `public/member` 以及是否允许 AI 读取。

## 5. 当前实现边界

当前已落地：

- `content` Studio service 层只允许 `admin`。
- 新建内容默认 `private + ai_access=denied`。
- 公开读取只返回 `published + public`。

当前不做：

- 不新增 `owner` role。
- 不设计多管理员层级、协作者、内容 ACL。
- 不把资源授权逻辑放到 `gateway`。
- 不实现 review TBAC 或 agent client 细粒度授权。
