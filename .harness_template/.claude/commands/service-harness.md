# Command: /service:harness

> **说明**：仅生成 Harness 工程制品（知识库目录 + 服务矩阵 + INDEX 更新），
> 适用于**用户已手动执行 `waterdrop new {service_name}` 后**，再调用本命令补全 Harness 接入。
>
> **触发**：用户输入 `/service:harness {service-name}` 或 `/service:harness`

---

## 与 /service:onboard 的区别

| 命令 | waterdrop new | Harness 制品 |
|-----|--------------|-------------|
| `/service:onboard` | ✅ 自动执行 | ✅ 自动生成 |
| `/service:harness` | ❌ 用户自己执行 | ✅ 自动生成 |

**使用场景**：
- 已经手动执行过 `waterdrop new chat`
- 或者业务代码目录是从其他方式创建的（非 waterdrop）
- 只需要补全 Harness 工程制品

---

## 执行步骤

### Step 1: 收集服务基本信息

如果用户没有提供，交互式询问（同 `/service:onboard`）：

```
请提供新服务的基本信息：

1. 服务名（英文小写，如 chat）：
2. 服务描述（中文，如 即时通讯服务）：
3. 通信协议：
   [1] HTTP（对外暴露，使用 grpc-gateway）
   [2] gRPC（内部服务）
4. 所属项目（如 api）：
5. 上游服务（哪些服务会调用本服务，逗号分隔，如 order,user）：
6. 下游服务（本服务会调用哪些服务，逗号分隔，如 user）：
7. 错误码段起始值（参考 context/team/error-code.md，如 5000）：
```

### Step 2: 验证业务代码目录存在

```
检查 api/{service-name}/ 目录是否存在...
```

- 目录**存在**：继续
- 目录**不存在**：提示用户先执行 `cd api && waterdrop new {service_name}`

### Step 3: 将服务加入工作区

将服务路径（如 `./api/{service-name}`）追加到仓库根目录的 `go.work` 文件中的 `use` 块内（如果尚未添加）。

### Step 4: 生成知识库目录

> ⚠️ **关键路径约束**：知识库目录必须创建在 `context/project/api/{service-name}/` 下，
> **绝对不能**创建在 `api/{service-name}/` 业务代码目录下。
> 业务代码目录 `api/{service-name}/` 中只放代码，不放任何 Harness 制品。

在 **`context/project/api/{service-name}/`** 下创建：

```
context/project/api/{service-name}/
├── INDEX.md
├── architecture.md
├── api-guide.md
├── operations-runbook.md
├── sop/.gitkeep
├── docs/api/.gitkeep
└── experience/.gitkeep
```

### Step 5: 更新服务矩阵

在 `.service-matrix/dependencies.yaml` 追加服务配置。

### Step 6: 更新项目 INDEX.md

在 `context/project/api/INDEX.md` 的服务列表表格中追加新服务行，并更新服务间依赖关系图。

### Step 7: 预留错误码段

在 `context/team/error-code.md` 中预留码段。

### Step 8: 输出完成摘要

```
✅ {service-name} 服务 Harness 制品已生成！

📁 新建文件：
  - context/project/api/{service-name}/INDEX.md
  - context/project/api/{service-name}/architecture.md
  - context/project/api/{service-name}/api-guide.md
  - context/project/api/{service-name}/operations-runbook.md
  - context/project/api/{service-name}/docs/api/.gitkeep

📝 更新文件：
  - go.work
  - .service-matrix/dependencies.yaml
  - context/project/api/INDEX.md
  - context/team/error-code.md

🔜 下一步：
  1. 根据通信协议清理 waterdrop 生成的默认模板代码（如去除不需要的 HTTP/gRPC 注册、删除 demo 文件等），并在清理后执行 `cd api/{service-name} && go build ./...` 验证编译通过
  2. 在 api/idl/{service-name}/ 下定义 proto 文件
  3. 在 context/project/api/{service-name}/INDEX.md 中补充关键约束
  4. 运行 /service:deps 验证服务依赖关系
```

---

## 示例

```
# 用户先手动执行 waterdrop
cd api && waterdrop new chat

# 然后告诉 Claude 生成 Harness 制品
/service:harness chat
```
