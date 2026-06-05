# 重新开始编码循环命令

## 描述
基于 `/coding:review` 的反馈，重新进入编码阶段，修正代码问题。

## 触发条件
用户输入 `/coding:restart`。

## 执行动作
1. 确定当前活动的 requirement id。**注意：如果用户输入的命令中自带了需求ID字段，则直接使用该ID，无需询问用户**。
2. 确认当前需求存在对应的 `gate-3-code-review.md`。如果不存在，提示先执行 `/coding:review`。
3. **读取门禁反馈**：读取 `gate-3-code-review.md` 中指出的代码规范问题、业务逻辑偏差、单测缺失等问题。
4. 展示需要修正的问题列表，引导用户逐一修正。
5. **【规范约束应用】**：在修正代码时，继续严格遵守 `context/team/development-sop.md` 和 `context/team/code-implementation-examples.md` 中的规范套路，严禁偏离约定。
6. **【DAO层实现约束】**：遵守原有 DAO 层约束，优先复用。
7. **【IDL 与编译约束】**：如果修正涉及 `.proto`，提示并自动执行 `waterdrop protoc --grpc --swagger *.proto`。
8. **【代码格式化约束】**：修正完成后，主动执行 `go fmt ./...`。
9. 提示用户修正完成后，重新运行 `/coding:review`。
