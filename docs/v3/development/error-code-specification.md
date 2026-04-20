# Beehive Blog v3 错误码规范

## 目标

建立 `v3` 全项目统一的业务错误码体系，确保：

- 业务错误码独立于 HTTP 状态码与 gRPC 状态码
- 客户端获得稳定、可检索、可文档化的错误语义
- 服务内部可以保留底层原因而不向客户端泄露敏感细节
- 后续 `content`、`search`、`realtime` 可直接复用同一套规则

## 设计原则

- 错误码使用六位整数业务码
- 业务错误码是领域真相源
- HTTP 与 gRPC 只负责 transport 适配
- 对外响应固定输出：
  - `code`
  - `message`
  - `reference`
  - `request_id`
- 对内日志可以记录底层 `cause`，但对外不能直接暴露 SQL、gRPC 原文、token、密码、client secret

## 编码规则

业务错误码采用三段式六位整数：

- 前两位：服务模块
- 中两位：错误类别
- 后两位：具体错误

### 服务模块段

- `10xxxx`：gateway
- `11xxxx`：identity
- `12xxxx`：content
- `13xxxx`：search
- `14xxxx`：realtime

### 错误类别段

- `xx01xx`：参数与请求错误
- `xx02xx`：认证错误
- `xx03xx`：授权错误
- `xx04xx`：资源状态错误 / 前置条件错误
- `xx05xx`：资源不存在 / 冲突
- `xx06xx`：依赖不可用 / 超时
- `xx99xx`：内部错误

## 当前已分配错误码

### Gateway

- `100101`：bad request
- `100201`：authorization required
- `100202`：invalid authorization scheme
- `100203`：access token invalid
- `100204`：access token inactive
- `100301`：access forbidden
- `100401`：gateway not ready
- `100601`：auth service unavailable
- `100602`：upstream timeout
- `109901`：internal error

### Identity

- `110101`：invalid argument
- `110201`：invalid credentials
- `110202`：invalid refresh token
- `110203`：refresh token expired
- `110204`：session revoked
- `110205`：account pending
- `110206`：account disabled
- `110207`：account locked
- `110401`：sso provider disabled
- `110402`：sso provider not ready
- `110403`：sso state invalid
- `110501`：resource not found
- `110502`：username already exists
- `110503`：email already exists
- `110601`：identity dependency unavailable
- `119901`：internal error

## message 规范

- `message` 面向客户端，必须是稳定英文短句
- `message` 不承载机器可读语义，机器可读语义由 `code` 承担
- 不允许把底层错误原文直接暴露给客户端

推荐示例：

- `invalid credentials`
- `access token is invalid`
- `service is not ready`
- `email already exists`

不推荐示例：

- `pq: duplicate key value violates unique constraint`
- `rpc error: code = Unavailable desc = connection refused`
- `token parse failed: signature is invalid`

## 代码使用要求

- 业务代码中不允许直接写裸整数错误码
- 必须通过 `pkg/errs` 中的常量名使用错误码
- 新增错误码时必须同步更新本规范文档
- 领域错误统一通过 `pkg/errs` 构造
- HTTP 错误输出统一通过 `pkg/errs/httpx`
- gRPC 错误输出统一通过 `pkg/errs/grpcx`
- 业务错误匹配首选 `errors.Is(err, errs.E(code))`
- 读取领域错误详情使用 `errs.Parse(err)` 或 `errors.As`
- `errors.Join` 只用于诊断聚合，不直接作为客户端主错误语义
- 禁止通过 `err.Error()`、`strings.Contains(err.Error(), ...)`、gRPC message、SQL message 做业务分支判断

## 新增错误码流程

1. 判断错误属于哪个服务模块
2. 判断错误属于哪个错误类别
3. 在对应段位中分配新的具体错误
4. 在 `pkg/errs/codes.go` 中新增常量
5. 在本文档中登记
6. 补充或更新对应测试
