# Command: /release:start

> **说明**：创建新版本的目录结构，以便管理此版本的需求、SQL变更、配置和脚本。
> **触发**：用户输入 `/release:start {version}` 或 `/release:start`

---

## 执行步骤

### Step 1: 收集版本号

如果用户没有提供版本号，询问：
- 请输入新版本的版本号（例如：v1.0.0）

### Step 2: 执行新建版本脚本

调用 Bash 工具执行：
```bash
bash scripts/new-release.sh {version}
```

### Step 3: 输出结果

输出提示：
```
✅ 成功创建版本目录：releases/{version}/

📝 下一步：
1. 您可以开始新建该版本下的需求，运行 `/requirement:start` 或 `/requirement:new`
2. 发布前需完善 releases/{version}/ 目录下的相关上线变更（SQL, 配置, 运维脚本）
```

---

## 参数

| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|-----|
| `version` | string | ❌ | 版本号，不填则交互式询问 |

---

## 示例

```
/release:start
/release:start v1.0.0
```
