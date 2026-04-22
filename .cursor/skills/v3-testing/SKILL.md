---
name: v3-testing
description: "Apply Beehive Blog v3 testing standards. Use when planning or writing tests and you need the repository's testing layers, Testcontainers defaults, and service/repo/logic test boundaries."
---

# v3 Testing

## Overview

Use this skill when a task touches tests.
在任务涉及测试时使用本 skill。

It covers test layering, Testcontainers defaults, and which layer should own which tests.
它负责测试分层、Testcontainers 默认方案，以及不同层的测试边界。

## Use When

- You are adding or modifying tests.
- You need to decide whether a case belongs in `service`, `repo`, or `logic`.
- You need to know the repository test defaults and fallback rules.

## Do

- Use the repository test layering model.
- Keep `service` as the main business test entry.
- Keep `logic` as smoke-test only.
- Use Testcontainers as the default integration-test path where the repository already requires it.
- When the task needs real gateway HTTP regression, chained auth validation, or repository-managed API test scripts, prefer the root `qa/` project.
- Treat `qa/` as the standard entry for Python + uv + pytest + locust based HTTP regression and load-test scaffolding.
- Use `qa/.env.example` and `qa/README.md` as the setup entrypoint for root-level QA runs.
- Run `uv sync --project qa` before the first local QA execution.
- Run `uv run --project qa python -m qa.scripts.check_env` before pytest or locust to confirm the target gateway is reachable and ready.
- Use `uv run --project qa pytest qa/tests` for repository-managed HTTP regression.
- Use `uv run --project qa locust -f qa/perf/locustfile.py` for the first-stage load-test skeleton.
- Keep HTTP chain validation, token propagation, and multi-endpoint regression flows inside `qa/flows/` and `qa/tests/`, not in ad-hoc external tools.

## Do Not

- Do not move heavy business coverage into `logic` tests.
- Do not default to local env fallback when the repository standard says Testcontainers first.
- Do not replace package-level Go tests with `qa/`; they serve different purposes.
- Do not put repository-owned HTTP regression cases back into Apifox or other external GUI tools.
- Do not edit generated or runtime-only artifacts to store QA state; keep test cases and scripts under `qa/`.
- Do not duplicate full coding or error/logging rules here.

## Hand-off

- For implementation-stage guidance, return to `$v3-during-task`.
- For final verification, continue with `$v3-finish-task`.

## References

- [测试规范](../../docs/v3/development/testing-conventions.md)
- [编码规范](../../docs/v3/development/coding-conventions.md)
- [QA 工程说明](../../qa/README.md)
