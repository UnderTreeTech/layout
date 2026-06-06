# Command: /release:incremental

> **说明**：增量发布梳理。针对某个小需求变更、BUGFIX 或 HOTFIX 进行临时发版，整理本次特定的变更（DDL、DML、配置、脚本），生成增量上线发布文档。
> **触发**：用户输入 `/release:incremental {version}` 或 `/release:incremental`

---

## 执行步骤

### Step 1: 确定发布版本及范围
1. 询问用户发版所在的版本号（例如 `v1.0.0`），如果命令参数中已提供则跳过。
2. 询问用户本次增量发布包含哪些具体的变更范围？（例如：“仅包含 T1234 的 bugfix”，“执行 user 表的 DDL”，“只需发布 order 服务配置”）。
3. 询问用户对这次增量发布的标签或别名（如 `hotfix-20260606-login-bug`）。

### Step 2: 过滤并收集变更产物
根据用户的范围约束，仅读取 `releases/{version}/` 下与本次增量相关的特定文件：
1. **涉及的需求/BugFix**：匹配对应的内容。
2. **涉及的数据库变更**：寻找本次必须执行的特定的 DDL/DML。
3. **涉及的配置与脚本变更**。

### Step 3: 生成增量发布文档
由 AI 根据收集的内容生成 `releases/{version}/RELEASE_PLAN_INCREMENTAL_{别名}.md` (增量上线发布文档)。
文档结构必须包含：
1. **增量背景**：为何临时发版，影响面（需指出涉及的具体服务）。
2. **变更详情**：仅列出本次抽取的 DDL、DML、配置与脚本变动，确保不会牵连其他未开发完的全量版本内容。
3. **执行流程**：步骤清晰（如先更新 SQL -> 再更新配置 -> 最后替换某一个微服务实例）。
4. **快速回滚方案**：只针对本次增量修改的紧急回滚措施。

### Step 4: 结果输出
输出成功提示，将生成的增量发布文档交给运维执行紧急插车或增量投产。

---

## 参数

| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|-----|
| `version` | string | ❌ | 版本号，不填则交互式询问 |
| `scope`   | string | ❌ | 本次增量发布包含的具体变更描述 |
| `alias`   | string | ❌ | 本次发布的别名（用于文件命名） |

## 示例

```
/release:incremental
/release:incremental v1.0.0
/release:incremental version=v1.0.0 alias=hotfix-order-bug scope=订单回调丢失问题
```
