---
name: code-review
description: "Multi-dimensional code review. Use when reviewing PRs, staged changes, or any code diff. Covers scope, dependencies, testing, performance, and security. Based on htmlpage.cn Cursor Code Review methodology."
---

# Code Review

## Overview

Multi-dimensional code review skill. Use this to systematically review code changes across five dimensions: scope, dependencies, testing, performance, and security. 用于对代码变更进行多维度系统化审查。

## When to Use

- Reviewing a PR or branch diff
- Pre-commit code self-review
- Reviewing changes before merging to main
- User asks to "review this code" or "code review"

## Five Review Dimensions

| 维度 | 重点问题 | 常见漏项 | 必须深查的场景 |
|------|---------|---------|-------------|
| **范围** | 改动是否越界，是否影响非目标模块 | 顺手重构、顺手改命名 | 多文件改动、共享组件 |
| **依赖** | 调用方、配置、路由、状态是否同步更新 | 只改定义不改使用方 | 接口、组件 props、store |
| **测试** | 是否补了单元/集成/冒烟验证 | 只说"建议测试"但不给清单 | 表单、支付、权限、发布 |
| **性能** | 首屏、体积、重复请求、阻塞资源 | 新增库体积、无缓存、重复渲染 | Landing 页、内容页、SSR 页面 |
| **安全** | 输入校验、权限边界、敏感信息泄漏 | 日志泄漏、token 暴露、越权访问 | 登录、用户资料、支付、上传 |

## Review Prompt Template

Use this prompt template when asking an AI to review code:

```
请只审查以下改动范围：
- 文件：<列出改动的文件>
- 目标：<一句话说明这次改动想解决什么>
- 不在范围内：<明确哪些模块不要展开>

请按以下维度输出审查结果：
1. 范围是否越界
2. 依赖与调用方是否漏改
3. 测试与回归项是否不足
4. 性能风险
5. 安全风险

输出要求：
- 只写有根据的问题，不要泛泛而谈
- 每条问题包含：严重级别、位置、原因、建议修复、回归检查点
- 如果没有发现问题，明确说明"未发现高置信问题"，并列出剩余验证风险
```

## Output Format

Each issue must include:

| Field | Description |
|-------|-------------|
| **严重级别** | 🔴 Critical / 🟠 High / 🟡 Medium / ⚪ Low |
| **位置** | File path and line range |
| **原因** | Why this is a problem |
| **建议修复** | Concrete fix suggestion |
| **回归检查点** | How to verify the fix is correct |

## Core Principles (5-Step Methodology)

1. **先限定范围** — define scope before reviewing
2. **再限定审查维度** — pick relevant dimensions from the five
3. **只接受高置信问题** — don't guess, flag only well-founded issues
4. **每条问题都带修复建议** — every issue must include a concrete fix
5. **最后转换成回归清单** — output actionable regression checkpoints

## Do

- Limit scope to the files/changes under review
- Check all five dimensions where applicable
- Mark issues with severity levels
- Provide concrete fix suggestions for every issue
- Include regression verification steps
- State "未发现高置信问题" when no issues found, with residual risk notes

## Do Not

- Do not review code outside the stated scope
- Do not flag style issues covered by linters/formatters
- Do not make vague suggestions without concrete fixes
- Do not skip dimensions without acknowledging them
- Do not review architectural decisions already approved
