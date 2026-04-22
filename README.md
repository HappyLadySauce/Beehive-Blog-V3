# Beehive Blog v3

Beehive Blog v3 是一个面向生产环境的博客系统后端工程，采用分层架构与契约优先（Contract-first）协作方式，聚焦可维护性、可扩展性与工程规范化。

## Features

- Contract-first API 设计（`v3/api` 与 `v3/proto`）
- 多服务模块化拆分（`services`）
- 通用能力沉淀（`pkg`）
- 质量保障与测试辅助（`qa`）
- CI 工作流支持（`.github/workflows`）

## Repository Structure

```text
.
├── data/              # 本地开发相关数据
├── docker/            # 容器与部署相关配置
├── docs/              # 项目文档
├── pkg/               # 跨服务可复用包
├── qa/                # 测试工具与客户端
├── scripts/           # 自动化脚本
├── services/          # 业务服务
├── sql/               # 数据库脚本
├── tools/             # 开发工具
└── v3/                # v3 合约定义与相关资产
```

## Quick Start

### Prerequisites

- Go（建议与 `go.mod` 版本保持一致）
- Docker / Docker Compose（可选，用于依赖服务）

### Install Dependencies

```bash
go mod download
```

### Run Tests

```bash
go test ./...
```

## Development Guidelines

- 贡献规范请参考 `CONTRIBUTING.md`
- 安全漏洞披露请参考 `SECURITY.md`
- 版本变更记录请参考 `CHANGELOG.md`

## License

本项目基于 [GNU GPLv3](./LICENSE) 开源发布。  
This project is licensed under GNU GPLv3.