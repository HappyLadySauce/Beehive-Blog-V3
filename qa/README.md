# Beehive Blog v3 QA 测试工程

本目录用于存放仓库内置的 HTTP 回归测试与性能测试工程。  
这套工程用于替代外部接口测试工具，把接口测试用例、测试脚本、认证串联逻辑和压测入口统一纳入仓库管理。

## 目标

- 对已经启动的服务执行真实的 gateway HTTP 回归测试。
- 把注册、登录、刷新、登出、鉴权访问等链路固化成可复用测试流程。
- 为后续压测保留稳定的 Python + Locust 入口。

## 目录说明

- `clients/`
  - HTTP API 客户端封装，统一处理 base URL、headers、timeout 和响应解析。
- `config/`
  - QA 环境配置模型与 `.env` 加载逻辑。
- `fixtures/`
  - 测试数据样例与唯一用户构造逻辑。
- `flows/`
  - 多步骤测试流程封装，例如注册并登录、刷新后再访问、登出后校验失效。
- `tests/`
  - pytest 回归测试用例。
- `perf/`
  - Locust 压测骨架与最小性能场景。
- `scripts/`
  - 辅助脚本，例如环境探测。

## 环境模型

第一阶段 QA 工程**不负责启动服务**。  
运行测试前，请确保以下组件已经启动：

- `docker/Infrastructure/docker-compose.yaml` 中的基础设施
- `identity` RPC 服务
- `gateway` HTTP 服务

默认的 gateway 基础地址是：

```text
http://127.0.0.1:8888
```

如果你的本地地址不同，请复制 `qa/.env.example` 为 `qa/.env` 后再修改配置。

## 安装依赖

在仓库根目录执行：

```powershell
uv sync --project qa
```

或者进入 `qa/` 目录后执行：

```powershell
uv sync
```

## 环境变量

推荐先复制模板文件：

```powershell
Copy-Item qa/.env.example qa/.env
```

常用变量包括：

- `BEEHIVE_QA_BASE_URL`
  - gateway 基础地址，默认 `http://127.0.0.1:8888`
- `BEEHIVE_QA_TIMEOUT_SECONDS`
  - HTTP 请求超时时间
- `BEEHIVE_QA_VERIFY_SSL`
  - 是否校验证书
- `BEEHIVE_QA_DEFAULT_PASSWORD`
  - 测试用户默认密码
- `BEEHIVE_QA_TEST_USERNAME_PREFIX`
  - 测试用户名统一前缀
- `BEEHIVE_QA_TEST_EMAIL_DOMAIN`
  - 测试邮箱域名
- `BEEHIVE_QA_ENABLE_SSO_TESTS`
  - 是否启用 SSO 回归测试

## 运行环境探测

在仓库根目录执行：

```powershell
uv run --project qa python -m qa.scripts.check_env
```

这个脚本会检查：

- `gateway` 基础地址是否可达
- `/healthz` 是否正常
- `/readyz` 是否正常

如果环境未就绪，它会直接给出失败信息，避免在 pytest 阶段才发现服务没有准备好。

## 运行 HTTP 回归测试

在仓库根目录执行全部回归：

```powershell
uv run --project qa pytest qa/tests
```

只跑认证相关测试：

```powershell
uv run --project qa pytest qa/tests/auth -q
```

只做测试收集，不真正执行：

```powershell
uv run --project qa pytest --collect-only qa/tests
```

## 运行 Locust 压测骨架

在仓库根目录执行：

```powershell
uv run --project qa locust -f qa/perf/locustfile.py
```

当前第一阶段已经提供的最小场景包括：

- `HealthzUser`
  - 持续访问 `/healthz`
- `AuthenticatedUser`
  - 注册测试用户后访问 `/api/v3/auth/me`

如果只想查看已注册的 Locust 用户类型：

```powershell
uv run --project qa locust -f qa/perf/locustfile.py --list
```

## 推荐工作流

推荐按下面顺序使用这套 QA 工程：

1. 启动基础设施
2. 启动 `identity`
3. 启动 `gateway`
4. 运行 `check_env`
5. 运行 pytest 回归
6. 需要时再运行 Locust

## 当前测试范围

第一阶段已经覆盖的 HTTP 回归范围：

- `GET /healthz`
- `GET /readyz`
- `POST /api/v3/auth/register`
- `POST /api/v3/auth/login`
- `POST /api/v3/auth/refresh`
- `POST /api/v3/auth/logout`
- `GET /api/v3/auth/me`

SSO 测试当前默认关闭，只有在显式开启相关环境变量时才会参与执行。

## 注意事项

- 脚本与测试场景输出的日志统一使用英文，便于后续自动化聚合。
- 断言同时覆盖 HTTP 行为和响应契约，不只检查状态码。
- 如果真实运行行为与 `v3/api/gateway.api` 不一致，应优先修正契约或服务实现，而不是静默放宽 QA 模型。
- 当前 QA 工程默认连接已启动服务，不会替你拉起 `gateway` 或 `identity`。
