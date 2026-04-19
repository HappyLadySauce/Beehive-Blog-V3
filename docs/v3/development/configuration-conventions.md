# Beehive Blog v3 配置规范

## 目标

统一 `v3` 的配置结构写法、go-zero `conf` tag 使用方式以及 `Validate()` 的职责边界。

## 基本原则

配置字段的校验职责固定分为两层：

### 第一层：tag 负责结构级约束

优先使用 go-zero `conf` tag 表达：

- `optional`
- `default=...`
- `range=...`
- `options=...`
- `env=...`

适合放进 tag 的内容：

- 默认值
- 数值范围
- 枚举值
- 可选性

### 第二层：`Validate()` 负责组合级与安全级约束

必须保留在 `Validate()` 的内容：

- 安全关键字段非空
- 跨字段关系校验
- 跨 provider 配置完整性校验
- URL、密钥、TTL 关系等组合逻辑

## tag 设计规则

### 什么时候用 `optional`

只用于真正允许缺省的字段，例如：

- Redis 用户名、密码
- TLS 开关
- 可选的 SSO `Scopes`
- 可选的回调基址

### 什么时候用 `default`

只要字段存在稳定且安全的默认值，就优先使用 `default`，例如：

- 端口
- 超时
- 连接池参数
- token TTL
- bcrypt cost

### 什么时候用 `range`

数值字段只要存在明确边界，就补 `range`，例如：

- 端口范围
- 超时秒数
- 连接池大小
- TTL 秒数
- 哈希成本

### 什么时候用 `options`

枚举型字符串配置必须补 `options`，例如：

- PostgreSQL `SSLMode`
- 日志模式
- 编码模式

### 什么时候不加 `optional`

基础连接和安全关键字段不加 `optional`，例如：

- `Postgres.Host`
- `Postgres.User`
- `Postgres.DBName`
- `StateRedis.Host`
- `Security.AccessTokenSecret`

## `Validate()` 设计规则

`Validate()` 只做 tag 做不了的事情，不重复做简单范围校验。

应重点承担：

- 密钥非空
- `RefreshTokenTTLSeconds > AccessTokenTTLSeconds`
- 启用的 SSO provider 必须具备：
  - `ClientID`
  - `ClientSecret`
  - `RedirectURL`
- URL 合法性校验

错误信息要求：

- 统一返回“配置路径 + 失败原因”
- 例如：
  - `Security.AccessTokenSecret 不能为空`
  - `SSO.GitHub.RedirectURL 必须包含合法的 scheme 和 host`

## Identity 配置示例

`services/identity/internal/config/config.go` 是当前 `v3` 的首个正式样板。

关键规则：

- `Postgres.Port` 使用 `default + range`
- `Postgres.SSLMode` 使用 `default + options`
- `StateRedis.Port` 使用 `default + range`
- `Security.AccessTokenSecret` 保持必填
- `Security.AccessTokenTTLSeconds`、`RefreshTokenTTLSeconds`、`StateTTLSeconds` 使用 `default + range`
- `SSO` 各 provider 只有在启用时，才由 `Validate()` 检查完整性

## 启动时序要求

所有服务统一遵循以下启动顺序：

1. 加载配置
2. 执行 `Validate()`
3. 初始化依赖
4. 启动服务

不允许把配置合法性判断散落到 `svc` 初始化细节中。
