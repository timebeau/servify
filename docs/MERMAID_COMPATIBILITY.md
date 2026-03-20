# Mermaid 兼容性

本仓库当前文档站使用 VuePress 2 默认主题和 Vite bundler。现阶段没有启用额外 Mermaid 渲染插件，因此 Mermaid 图需要遵守以下约束。

## 当前约束

- Markdown 中不要默认假设 ```mermaid 会被前端自动渲染
- 需要长期保真保存的图，应同时提供：
  - Mermaid 源码块
  - 对应的文字版步骤说明
- 评审或 PR 中新增 Mermaid 图时，必须先本地执行 `npm -C docs run build`
- 如果图是关键交付内容，优先附带导出的 PNG/SVG 到 `docs/public/`

## 推荐写法

1. 先写纯文本的小节说明流程、节点含义和关键分支
2. Mermaid 仅作为辅助表达，不承担唯一信息源
3. 控制图规模，避免超宽布局和过深嵌套
4. 节点文案使用短句，避免中英文混排过长导致换行不可控

## 后续演进

- 如果后续正式引入 Mermaid 插件，需要补：
  - VuePress 插件版本锁定
  - dark/light theme 可读性校验
  - CI 构建快照或最小渲染冒烟用例
