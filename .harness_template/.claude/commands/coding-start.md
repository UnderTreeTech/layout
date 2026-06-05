# 开始编码循环命令

## 描述
进入 Vibe Coding 阶段，基于拆解的任务开始编写代码。

## 触发条件
用户输入 `/coding:start`。

## 执行动作
1. 确定当前活动的 requirement id。**注意：如果用户输入的命令中自带了需求ID字段，则直接使用该ID，无需询问用户**。
2. 确认当前需求已通过设计门禁（检查 `gate-2-design-review.md`）。
3. 读取 `tasks/features.json`，展示待开发的任务列表。
4. 引导用户选择当前要开发的任务。
5. **【规范加载】**：必须强制加载并读取 `context/team/development-sop.md` 以及标准样板库 `context/team/code-implementation-examples.md`。
6. **【编码约束】**：准备开始 AI 辅助编码时，必须严格按照 SOP 中规定的套路（HTTP、gRPC 或普通方法的标准实现路径）来生成代码。**严禁 AI 自行发挥或偏离团队框架约定**。
7. **【DAO层实现约束】**：正常情况下，底层的 `dao`、`model` 及 `iface` 代码用户会提前通过 `xo` 组件自动生成好，AI 在编码阶段应优先直接引用和使用现有代码；**仅当文件不存在时**，询问用户是否自己去生成代码，如果用户拒绝才允许 AI 严格按照 `code-implementation-examples.md` 中的样板代码去手动生成或补齐这三层代码。
8. **【Proto 编译约束】**：如果编码过程中涉及修改 `.proto` 文件，必须主动执行命令 `waterdrop protoc --grpc --swagger *.proto`（在对应的 `api/idl/{service_name}` 目录下）生成新的 stub 文件。
9. **【代码格式化约束】**：每次完成一定阶段的代码编写后，必须主动执行 `go fmt ./...` 及 `go build ./...` 命令来格式化及校验当前修改的服务代码。
10. **【HTTP API接口文档】**：HTTP接口的设计必须严格遵循 `context/team/api-desing.md` 规范，并将接口文档生成至 `context/project/api/{service_name}/docs/api/` 目录下。