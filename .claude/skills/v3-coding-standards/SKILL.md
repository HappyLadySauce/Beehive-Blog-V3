---
name: v3-coding-standards
description: "Apply Beehive Blog v3 coding standards. Use when handwritten code is being designed or edited and you need the repository's layer boundaries, directory responsibilities, config rules, and comment conventions."
---

# v3 Coding Standards

## Overview

Use this skill for general implementation standards.
在需要遵守项目通用编码规范时使用本 skill。

It covers repository structure, layer boundaries, config boundaries, and comment conventions.
它负责仓库结构、分层边界、配置边界和注释规范。

## Use When

- You are editing handwritten business code.
- You need to decide which layer owns a change.
- You need to apply repository-wide config or comment conventions.

## Do

- Keep implementation inside the correct directory boundary.
- Follow the repository service layering model.
- Use bilingual comments with English above Chinese.
- Follow the config and validation conventions from the development docs.

## Do Not

- Do not duplicate error/logging detail here.
- Do not inline test strategy detail here.
- Do not treat `gateway` as a business orchestration layer.

## Hand-off

- For error/log changes, continue with `$v3-error-and-logging`.
- For test work, continue with `$v3-testing`.
- For final checks, continue with `$v3-finish-task`.

## References

- [编码规范](../../docs/v3/development/coding-conventions.md)
- [配置规范](../../docs/v3/development/configuration-conventions.md)
