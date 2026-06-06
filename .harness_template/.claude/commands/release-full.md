# Command: /release:full

> **说明**：全量发布梳理。整理目标版本下所有需求的变更（DDL、DML、配置、脚本），生成全量上线发布及运维执行文档。
> **触发**：用户输入 `/release:full {version}` 或 `/release:full`

---

## 执行步骤

### Step 1: 确定发布版本
询问用户准备全量发布的版本号（例如 `v1.0.0`），如果命令参数中已提供则跳过。确认 `releases/{version}/` 目录是否存在，若不存在提示用户检查。

### Step 2: 收集全量变更产物
使用工具读取并汇总该版本下的所有变更文件：
1. **需求列表**：遍历 `releases/{version}/requirements/` 下的需求，提取准备发布的需求及 BugFix 列表。
2. **数据库变更**：读取 `releases/{version}/sql/ddl/` 和 `releases/{version}/sql/dml/` 下的所有 SQL 文件，按序号升序排列。
3. **配置变更**：读取 `releases/{version}/configs/diff/` 下的配置变更/比对文件。
4. **脚本变更**：读取 `releases/{version}/scripts/pre-deploy/` 和 `post-deploy/` 下的运维脚本。
5. **回滚变更**：梳理对应的 `rollback` 策略。

### Step 3: 生成全量发布文档
根据上述收集的信息，由 AI 自动编写并合并成 `releases/{version}/RELEASE_PLAN_FULL.md` (全量上线发布文档)。
文档结构必须包含：
1. **版本信息**：版本号、计划发布时间、整体影响面概述。
2. **需求集**：本次全量发布包含的所有功能与缺陷修复清单。
3. **执行顺序**：
   - 部署前准备（执行 `pre-deploy` 脚本等）。
   - 数据库执行计划（详细列出将要执行的 DDL 和 DML SQL 文件名及目的）。
   - 服务代码更新（后端服务镜像替换/重启逻辑）。
   - 增量配置推送（`configs` 变更项）。
   - 部署后操作（执行 `post-deploy` 脚本、缓存预热等）。
4. **验证与验收**：上线后的冒烟测试核心验证点。
5. **回滚方案**：如果整体发版失败，相关的 `rollback` 脚本与数据库回滚指南。

### Step 4: 结果输出
输出成功提示，引导用户将生成的 `RELEASE_PLAN_FULL.md` 分发给运维团队，或将其转录至自动化发布的流水线系统中。

---

## 参数

| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|-----|
| `version` | string | ❌ | 版本号，不填则交互式询问 |

## 示例

```
/release:full
/release:full v1.0.0
```
