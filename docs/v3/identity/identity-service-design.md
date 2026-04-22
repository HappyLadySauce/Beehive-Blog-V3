# Beehive Blog v3 Identity 服务设计

## 1. 目标

本文件用于定义 `v3` 阶段 `identity` 服务的职责边界、与 `gateway` 的关系、与业务服务的协作方式，以及第一阶段的实现收口。

核心目标：

- 明确 `identity` 到底负责什么
- 明确 `gateway` 与 `identity` 如何协作
- 明确用户、认证、会话、SSO 在第一阶段如何落地
- 避免后续在 `gateway`、`identity`、未来 `user-service` 之间出现边界漂移

## 2. 设计原则

### 2.1 identity 是认证与基础账户服务

`identity` 负责：

- 本地账号注册与登录
- SSO 登录接入
- access token / refresh token
- 当前用户身份解析
- 会话管理与吊销
- 用户基础账户信息
- 用户角色与账号状态

### 2.2 gateway 只做认证前置，不做授权裁决

`gateway` 负责：

- 接收外部 HTTP / WS 请求
- 提取 access token
- 调用 `identity` 完成登录、刷新、当前用户、token 校验
- 把认证后的身份上下文透传给下游服务

`gateway` 不负责：

- 资源级授权
- 用户资料真相
- 业务编排

### 2.3 授权下沉到业务服务

`identity` 提供“你是谁”的真相。

业务服务负责判断：

- 你能不能读某条资源
- 你能不能修改某条资源
- 你是不是该资源所有者
- 你是否具备 `admin` 管理能力

### 2.4 第一阶段不拆 user-service

当前阶段不单独拆分 `user-service`。

基础用户资料与账户能力统一归 `identity`。

### 2.5 认证主形态为本地账号 + SSO 并存

第一阶段同时支持：

- 本地账号：`username 或 email + password`
- SSO：联邦身份登录

这意味着：

- `identity` 必须支持统一 `user` 主体
- 本地凭证和联邦身份都绑定到同一个 `user`
- SSO 不是未来注释位，而是一等设计对象

### 2.6 第三阶段开放 GitHub / QQ / WeChat SSO

第三阶段的实际交付范围固定为：

- `GitHub`：完整开放 `StartSsoLogin + FinishSsoLogin`
- `QQ`：完整开放 `StartSsoLogin + FinishSsoLogin`
- `WeChat`：完整开放 `StartSsoLogin + FinishSsoLogin`

补充说明：

- QQ 采用 Web OAuth 流程，以 `openid` 作为主体标识
- WeChat 采用网站/扫码 OAuth 流程，优先使用 `unionid`，没有时回退 `openid`
- 实现继续复用统一的 `StartSsoLogin / FinishSsoLogin` 抽象，不新增 provider 专属 RPC

## 3. 服务职责边界

### 3.0 内部实现层次

`identity` 当前正式实现层次为：

```text
server -> logic -> service -> repo -> entity
```

说明：

- `server` / `logic`：只做 gRPC transport 适配
- `service`：负责核心用例编排、事务边界、审计与认证流程
- `repo`：负责持久化访问
- `entity`：负责 GORM 表结构映射

### 3.1 identity 负责

- 本地账号注册
- 本地账号登录
- GitHub / QQ / WeChat SSO 登录接入
- access token 签发
- refresh token 签发与轮换
- token introspection
- 当前用户身份信息查询
- 会话管理
- 单会话吊销与登出
- 用户角色与账号状态维护
- 身份审计日志

### 3.2 identity 不负责

- 内容、评论、搜索等业务主数据
- 资源级权限矩阵裁决
- 内容可见性、状态、AI 访问控制
- 网关层路由
- 无主语聚合接口

### 3.3 与 gateway 的协作关系

普通认证链路：

```text
Client -> Gateway -> Identity -> Gateway -> Client
```

受保护接口链路：

```text
Client -> Gateway -> Identity(Introspect) -> Gateway -> Business Service
```

刷新链路：

```text
Client -> Gateway -> Identity(Refresh) -> Gateway -> Client
```

### 3.4 与业务服务的协作关系

业务服务只依赖 `identity` 提供的身份上下文，不依赖 `gateway` 做资源授权。

推荐下游上下文字段至少包括：

- `user_id`
- `role`
- `account_status`
- `session_id`
- `auth_source`

## 4. 第一阶段用户与权限收口

### 4.1 角色模型

第一阶段统一采用：

- `guest`
- `member`
- `admin`

说明：

- `guest` 为匿名访客，不落库为业务用户
- `member` 为默认注册用户
- `admin` 为平台管理角色

### 4.2 账号状态模型

第一阶段统一采用：

- `pending`
- `active`
- `disabled`
- `locked`

默认值建议：

- 若第一阶段不引入注册审核，则默认新账号状态为 `active`

### 4.3 权限边界

`identity` 负责：

- 主体身份是否存在
- 主体角色是什么
- 主体账号是否可用

业务服务负责：

- 主体是否能访问资源
- 主体是否能执行具体动作
- 主体是否满足资源所有权或管理权限条件

## 5. Token 与会话策略

### 5.1 access token

采用短期 JWT。

主要用途：

- 给 `gateway` 做请求认证
- 给下游服务承载身份上下文

最少包含：

- `sub` / `user_id`
- `role`
- `account_status`
- `session_id`
- `auth_source`
- `exp`

### 5.2 refresh token

采用服务端持久化 refresh token。

特性要求：

- 支持轮换
- 支持失效
- 支持按会话吊销
- 支持多端并存

### 5.3 多会话模型

第一阶段采用多会话设计：

- 一个用户可以同时登录多个设备/客户端
- 每次登录创建一个 `user_session`
- refresh token 绑定到会话
- 登出默认按当前会话吊销

### 5.4 安全要求

第一阶段必须明确：

- 密码使用成熟哈希算法存储
- refresh token 不明文持久化
- token 失效与账号状态变更可联动
- `disabled` / `locked` 用户即使持有旧 token，也应在 introspection 时被拦截

## 6. SSO 设计原则

SSO 采用统一联邦身份抽象，不在 proto 层绑定单一厂商。

当前第三阶段的实现状态为：

- 已完整落地：`GitHub`
- 已完整落地：`QQ`
- 已完整落地：`WeChat`

约束：

- 一个联邦身份只能绑定一个 `user`
- 一个 `user` 可绑定多个联邦身份
- SSO 登录成功后仍应生成平台自己的会话与 token

## 7. 当前结论

`v3` 第一阶段的 `identity` 已收口为：

**一个同时负责本地账号、GitHub / QQ / WeChat SSO、会话、token、基础账户信息的认证与基础账户服务。**

补充说明：

- `QQ` 使用 `openid` 作为联邦主体标识
- `WeChat` 优先使用 `unionid`，没有时回退 `openid`
- `gateway` 继续只做 provider 透传和 transport 适配，不增加 provider 业务分支
