# 需求重写执行命令

> **说明**：基于 `/requirement:review` 的反馈，重新拆解并更新需求。
> **前置条件**：已通过 `/requirement:write` 生成过需求文档，并由 `/requirement:review` 给出了门禁反馈。
> **触发**：用户输入 `/requirement:rewrite`

---

## 执行步骤

### Step 1: 确定当前需求 ID 和文档路径
1. 加载当前活动的 requirement id（可从 `.harness/local.yaml` 或询问用户获得）。**注意：如果用户输入的命令中自带了需求ID字段（如 `/requirement:rewrite T123`），则直接使用该ID，无需询问用户**。
2. 定位到对应的 `requirements/{requirement-id}/requirement.md` 和 `gate-1-requirement-review.md`。如果 review 文件不存在，提示用户先执行 `/requirement:review`。

### Step 2: 收集门禁反馈
读取 `gate-1-requirement-review.md` 中指出的问题点和修改建议。如果用户在命令后附加了额外的修改意见文本，也一并收集。

### Step 3: 读取需求模板
读取框架级模板 `context/harness-framework/templates/requirement-template.md` 的规范结构。

### Step 4: AI 修正与重新组装
根据原需求文档、门禁反馈、用户补充意见和模板要求，针对性地修正以下内容：
1. **背景和目标**：修正不清晰的背景和难以度量的目标。
2. **非目标**：补充遗漏的非目标。
3. **用户故事**：规范格式，确保描述清晰。
4. **功能需求条目**：细化或拆分过大的需求条目，补充缺失的验收标准。
5. **非功能需求**：完善性能、可用性、安全等要求。
6. **影响面分析**：根据反馈修正影响面分析。

### Step 5: 更新需求文档
将修正后的 Markdown 内容写回 `requirements/{requirement-id}/requirement.md` 文件。

### Step 6: 引导用户下一步
输出提示：
```
✅ 需求文档已基于评审意见修正并重新写入：requirements/{requirement-id}/requirement.md

📝 下一步：
1. 请人工走查确认修正后的需求文档。
2. 确认无误后，重新运行 /requirement:review 进行需求评审门禁检查。
```
