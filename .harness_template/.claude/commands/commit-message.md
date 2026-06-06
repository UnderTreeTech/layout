# 生成 Commit Message 命令

## 描述
分析当前代码库中已暂存（staged）、未暂存（unstaged）或最近的修改记录，自动生成符合约定式提交规范（Conventional Commits）的中英文 Commit Message 供用户参考使用。

## 触发条件
用户输入 `/commit-message` 或相关的自然语言（如：“给我一个 commit message”、“总结一下最近的代码修复”等）。

## 执行动作
1. **收集变更上下文**：
   - 优先通过 `git diff --cached` 获取已暂存的更改。
   - 如果没有暂存更改，则通过 `git diff` 获取当前工作区未暂存的更改。
   - 也可以视情况结合 `git status` 了解新文件/删除文件情况。
2. **分析变更内容**：
   - 归纳变更的核心范围（Scope），如 API、Team SOP、配置、前端等。
   - 判断变更的类型（Type），如 `feat`（新功能）、`fix`（修复）、`docs`（文档）、`style`（格式）、`refactor`（重构）、`chore`（构建过程或辅助工具变动）等。
3. **生成总结内容**：
   - 生成符合 [Conventional Commits](https://www.conventionalcommits.org/) 规范的 Header 和 Body。
   - Body 中使用要点（bullet points）清晰描述修改了什么，以及为什么修改。
4. **输出格式要求**：
   - 必须同时输出**英文版本**和**中文版本**。
   - 提供直接可以复制执行的 `git commit -m "..." -m "..."` 命令行代码块，方便用户直接复制粘贴到终端中运行。

### 输出示例模版

```markdown
这是为您生成的 Commit Message：

#### 英文版 (English Version)
\`\`\`text
type(scope): short description

- Detail 1
- Detail 2
\`\`\`

可以直接执行的命令：
\`\`\`bash
git commit -m "type(scope): short description" -m "- Detail 1
- Detail 2"
\`\`\`

---

#### 中文版 (Chinese Version)
\`\`\`text
type(scope): 简短描述

- 详情 1
- 详情 2
\`\`\`

可以直接执行的命令：
\`\`\`bash
git commit -m "type(scope): 简短描述" -m "- 详情 1
- 详情 2"
\`\`\`
```