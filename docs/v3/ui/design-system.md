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

首批自研组件：

- `BaseButton`：按钮状态、尺寸、loading。
- `BaseInput`：表单输入、错误文案。
- `BaseCard`：通用面板容器。
- `BaseBadge`：类型、状态、标签。
- `BaseTabs`：轻量切换器。
- `EmptyState`：空状态。
- `PublicLayout`：公开站布局。
- `StudioLayout`：Studio 工作台布局。

图标优先使用 `lucide-vue-next`，避免手写 SVG 图标。
