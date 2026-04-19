# Beehive Blog v2 服务启动顺序与落地步骤

## 1. 目标

本文件定义 v2 从当前仓库状态继续推进的顺序。

当前不是“从零写文档”，而是“在已重组的仓库上继续补服务代码”。

## 2. 当前已完成

已完成项：

- 仓库已切换为 `api/ + proto/ + services/ + pkg/`
- `domain/` / `repository/` 预留已取消
- `gateway` 已用 `goctl api go` 生成
- `identity.proto`
- `content.proto`
- `search.proto`
- `identity/content/search` 目录骨架已就位

## 3. 当前第一优先级

当前第一优先级不是继续扩文档，而是把第一批服务真正跑起来。

顺序固定为：

1. `identity`
2. `content`
3. `gateway`
4. `search`
5. `indexer`

## 4. Phase 1：补齐基础依赖

目标：

- 确保 RPC 可以正式生成

动作：

1. 安装 `protoc`
2. 保留已安装的 `goctl`
3. 校验 `goctl --version`
4. 校验 `protoc --version`

完成标准：

- 可以执行 `goctl rpc protoc`

## 5. Phase 2：生成 `identity`

目标：

- 生成第一个 RPC 服务

动作：

1. 固化 `proto/identity.proto`
2. 用 `goctl rpc protoc proto/identity.proto --go_out=. --go-grpc_out=. --zrpc_out=services/identity`
3. 实现最小能力：
   - `Register`
   - `Login`
   - `RefreshToken`
   - `GetCurrentUser`

完成标准：

- `identity` 服务可启动

## 6. Phase 3：生成 `content`

目标：

- 生成核心内容 RPC 服务

动作：

1. 固化 `proto/content.proto`
2. 用 `goctl rpc protoc ...` 生成 `services/content`
3. 实现最小能力：
   - `CreateContent`
   - `UpdateContent`
   - `GetContent`
   - `ListContents`
   - `UpdateContentStatus`

完成标准：

- owner 可通过 RPC 维护内容

## 7. Phase 4：接通 `gateway`

目标：

- 让 `gateway` 真正转发到内部服务

动作：

1. 在 `services/gateway` 增加 RPC client 配置
2. 把 auth 路由接到 `identity`
3. 把内容路由接到 `content`
4. 补统一响应与错误码映射

完成标准：

- 前端只访问 `gateway`

## 8. Phase 5：生成 `search`

目标：

- 让站内搜索具备最小可用能力

动作：

1. 固化 `proto/search.proto`
2. 用 `goctl rpc protoc ...` 生成 `services/search`
3. 先接 Meilisearch
4. 实现：
   - `Search`
   - `Suggest`
   - `Related`

完成标准：

- 公开内容可搜索

## 9. Phase 6：补 `indexer`

目标：

- 自动更新搜索索引

动作：

1. 固化 worker 配置
2. 接内容变更事件或事件表
3. 实现 reindex

完成标准：

- 内容变更后索引能自动同步

## 10. 第一批必须打通的接口

### `identity`

- `Register`
- `Login`
- `RefreshToken`
- `GetCurrentUser`

### `content`

- `CreateContent`
- `UpdateContent`
- `GetContent`
- `ListContents`
- `UpdateContentStatus`

### `search`

- `Search`

### `gateway`

- `/api/v2/auth/*`
- `/api/v2/public/articles`
- `/api/v2/studio/contents*`
- `/api/v2/search/query`

## 11. 当前结论

当前最合理的推进顺序是：

**先补 `protoc`，再生成 RPC 服务，最后把 `gateway` 接成唯一对外入口。**
