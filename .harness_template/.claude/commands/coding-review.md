# 代码审查门禁执行命令

## 描述
执行编码阶段完成后的门禁检查（即代码审查）。

## 触发条件
用户输入 `/coding:review`。

## 执行动作
1. 加载当前活动的 requirement id。**注意：如果用户输入的命令中自带了需求ID字段，则直接使用该ID，无需询问用户**。
2. 收集当前工作区未提交的代码变更（diff）。
3. 调用 `code-review-report` Skill 或相关审查 Agent。
4. 校验代码是否符合团队规范、业务逻辑是否符合设计文档、核心逻辑是否有单测覆盖。
5. **IDL 定义与生成**：IDL定义在 `api/idl/{service_name}/`，务必校验IDL是否有正确生成，如未生成报错并反馈需执行 `waterdrop protoc --grpc --swagger *.proto` 来生成。
6. 在api目录下执行 `go fmt ./...`，`go vet ./...` 及 `go build ./...` 做格式化统一及验证编译是否能通过（api/idl目录除外）。
7. 生成审查报告并**以追加（Append）的方式**写入 `requirements/{version}/{requirement-id}/gate-3-code-review.md`。每次写入时需带有当前时间戳或审查轮次标识，以留存历次代码审查的历史记录。同时，在报告中必须专门记录**“本次踩过的坑/发现的代码缺陷”**，以便追溯和总结编码阶段踩过的坑。
8. 若通过，则提示用户可以提交代码；若不通过，则指出需要修改的地方。