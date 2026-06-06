# Command: /requirement:start

> **说明**：一步完成需求的新建与撰写，直接从零开始到生成完整需求文档。
> **触发**：用户输入 `/requirement:start` 或 `/requirement:start {需求标题}`

---

## 执行步骤

### Step 1: 收集基本信息与需求输入

如果用户没有提供需求标题，询问：
- 需求标题
- 需求 ID（如有 Geelib/TAPD/Jira ID）
- 优先级（P0/P1/P2/P3）
- 产品需求输入：提供产品需求文档的 docId，或者直接输入一段产品需求描述。（如果提供 docId，请使用 knowledge-skill 的 cli 能力读取需求内容，执行命令为：`editor-cli read <docId>`）

### Step 2: 加载上下文

按顺序加载：
1. `context/team/INDEX.md` - 团队规范
2. `context/harness-framework/main-process-numbering.md` - 流程规范
3. `.service-matrix/dependencies.yaml` - 服务矩阵
4. `context/project/{project-name}/INDEX.md` - 项目知识
5. `context/harness-framework/templates/requirement-template.md` - 需求模板

### Step 3: 创建目录结构

```
requirements/{version}/{requirement-id}/
├── requirement.md          # 需求文档（后续将写入 AI 生成的内容）
├── status.json             # 初始化为 INIT 阶段
└── tasks/
    └── .gitkeep
```

### Step 4: AI 拆解与组装需求文档

根据用户输入的产品文档内容和读取到的需求模板，执行如下拆分与补齐动作：
1. **背景和目标**：从产品文档中提炼为什么要做，以及可度量的目标。
2. **非目标**：识别当前范围不包含的内容。
3. **用户故事**：将业务描述转化为“作为...希望...以便...”的格式。
4. **功能需求条目**：将产品模块拆分成原子级的需求条目（REQ-xxx），并提取/生成具体可测试的验收标准。
5. **非功能需求**：识别性能、可用性、安全等要求。
6. **影响面分析**：从服务矩阵中初步提取可能涉及的模块/服务，并留空待开发者进一步补充确认。
7. **忽略项**：无须关心一切与UI相关的设计，不用体现出来。

生成最终 Markdown 内容，并写入 `requirements/{version}/{requirement-id}/requirement.md` 文件。

### Step 5: 初始化 status.json

```json
{
    "requirement_id": "{requirement-id}",
    "title": "{需求标题}",
    "current_stage": "REQUIREMENT_DEFINING",
    "current_step": "2.1",
    "gates_passed": [],
    "assignee": "",
    "created_at": "{当前时间}",
    "last_updated": "{当前时间}",
    "branch": "feature/{devops-name}/{geelib-id}"
}
```

### Step 6: 引导用户下一步

输出提示：
```
✅ 需求目录已创建且文档已生成写入：requirements/{version}/{requirement-id}/requirement.md

📝 下一步：
1. 请人工走查并微调生成的需求文档，特别是"影响面分析"和"验收标准"部分。
2. 确认无误后，运行 `/requirement:review` 进行需求评审门禁检查。
3. 通过后运行 `/design:new` 进入设计阶段。

💡 提示：
- 影响面分析请参考 .service-matrix/dependencies.yaml
- 需求模板说明请参考 context/harness-framework/templates/requirement-template.md
```

---

## 参数

| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|-----|
| `title` | string | ❌ | 需求标题，不填则交互式询问 |
| `id` | string | ❌ | 需求 ID，不填则自动生成 |
| `docId` | string | ❌ | 产品需求文档 ID，不填则交互式询问 |

---

## 示例

```
/requirement:start
/requirement:start 用户续费功能优化
/requirement:start id=T12345678 docId=doc_12345 用户续费功能优化
```
