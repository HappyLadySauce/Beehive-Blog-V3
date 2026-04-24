# UI Design System

## 1. 设计取向

首版采用分区混合设计：

- Public Web：参考 AnZhiYu 的阅读结构与 Notion 的温和留白，强调内容可读性。
- Studio：参考 Linear/Airtable 的高密度工作台结构，强调扫描、比较和重复操作。
- 文档：参考 Mintlify 的清晰分层，以工程可读性为优先。

不直接复制任一参考站点，而是抽象出 Beehive 自己的设计令牌与组件规则。

## 2. 设计令牌

核心 token 在 `ui/src/app/styles/app.css` 与 `ui/uno.config.ts` 中维护：

- `--bb-color-ink`：主文字
- `--bb-color-muted`：次级文字
- `--bb-color-paper`：页面背景
- `--bb-color-surface`：面板背景
- `--bb-color-line`：边框与分割线
- `--bb-color-honey`：创作/重点状态
- `--bb-color-leaf`：成功/知识沉淀状态
- `--bb-color-blue`：链接/导航状态
- `--bb-color-violet`：Studio/AI 状态

组件圆角默认不超过 `8px`。文字不使用视口宽度缩放，桌面和移动端通过断点调整布局密度。

## 3. 响应式规则

断点固定为：

| Token | Width |
| --- | --- |
| `sm` | `640px` |
| `md` | `768px` |
| `lg` | `1024px` |
| `xl` | `1280px` |
| `2xl` | `1536px` |

布局规则：

- Public Web：移动端单列，桌面端使用内容网格和 sticky 顶栏。
- Studio：桌面端侧栏 + 顶栏，移动端侧栏折叠为抽屉导航。
- 固定格式 UI 需要稳定尺寸，避免 hover、加载态和标签导致布局跳动。
- 所有页面必须避免横向滚动和文本溢出。

## 4. 基础组件

首批自研组件按用途分组：

- 操作：
  - `BaseButton`：文本按钮、提交按钮、loading。
  - `IconButton`：仅图标操作，必须提供 `label` 作为可访问名称。
- 表单：
  - `BaseInput`：单行输入。
  - `BaseTextarea`：多行输入。
  - `BaseSelect`：小型选项集。
  - `BaseCheckbox`：多选或确认项。
  - `BaseSwitch`：二元开关。
- 信息结构：
  - `BaseCard`：通用面板容器。
  - `BaseBadge`：类型、状态、标签。
  - `BaseTabs`：轻量切换器。
  - `PageHeader`：页面标题与主操作区。
  - `SectionHeader`：模块标题与局部操作区。
  - `DataTable`：Studio 轻量数据表。
- 反馈：
  - `StatusAlert`：信息、成功、警告、危险提示。
  - `EmptyState`：空状态。
  - `LoadingSkeleton`：加载占位。
- 布局：
  - `PublicLayout`：公开站布局。
  - `StudioLayout`：Studio 工作台布局。

使用边界：

- 页面优先组合 `PageHeader`、`SectionHeader`、`BaseCard` 和业务组件，不直接堆大量裸样式。
- Studio 列表优先使用 `DataTable`；卡片网格仅用于概览、看板或 Public Web 内容预览。
- 表单控件必须通过 `v-model` 暴露值变化，不在组件内部保存业务状态。
- 反馈组件只展示状态，不负责发起请求、重试或导航。

图标优先使用 `lucide-vue-next`，避免手写 SVG 图标。
