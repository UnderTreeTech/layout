# 快速编码启动命令

## 描述
用于处理临时提出的小需求或 BUG 修复，跳过标准的需求评审与设计门禁等繁杂检查，直接进入 Vibe Coding 阶段开始编写代码。但全程开发流程与代码规范仍需严格遵循团队标准。

## 触发条件
用户输入 `/coding:quick-start` 或 `/coding:hotfix`。

## 执行动作
1. **明确任务目标**：确认用户当前要实现的小需求或要修复的 BUG 的具体内容（如果用户在触发命令时已经提供则直接使用，无需再问）。
2. **跳过门禁检查**：无需确认 requirement id，无需检查 `gate-2-design-review.md` 及解析 `tasks/features.json`。
3. **【规范加载】**：必须强制加载并读取 `context/team/development-sop.md` 以及标准样板库 `context/team/code-implementation-examples.md`。
4. **【编码约束】**：准备开始 AI 辅助编码时，必须严格按照 SOP 中规定的套路（HTTP、gRPC 或普通方法的标准实现路径）来生成代码。**严禁 AI 自行发挥或偏离团队框架约定**。
5. **【DAO层实现约束】**：底层的 `dao`、`model` 及 `iface` 代码必须由用户通过 `xo` 组件自动生成。AI 在编码阶段应优先直接引用和使用现有代码。**绝对禁止 AI 手写生成这三层的实现代码**。如果缺失相关代码，必须提示用户通过命令行工具生成。
6. **【Proto 编译约束】**：如果编码过程中涉及修改 `.proto` 文件，必须主动执行命令 `waterdrop protoc --grpc --swagger *.proto`（在对应的 `api/idl/{service_name}` 目录下）生成新的 stub 文件。
7. **【代码格式化约束】**：每次完成一定阶段的代码编写后，必须主动执行 `go fmt ./...` 及 `go build ./...` 命令来格式化及校验当前修改的服务代码。
8. **【HTTP API接口文档】**：HTTP接口的设计必须严格遵循 `context/team/api-desing.md` 规范，并将接口文档生成至 `context/project/api/{service_name}/docs/api/` 目录下。