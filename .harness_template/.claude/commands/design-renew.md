# 重新设计方案命令

## 描述
基于 `/design:review` 的反馈，重新修正方案设计和任务拆解。

## 触发条件
用户输入 `/design:renew`。

## 执行动作
1. 确定当前活动的 requirement id。**注意：如果用户输入的命令中自带了需求ID字段，则直接使用该ID，无需询问用户**。
2. 确认当前需求存在对应的 `design.md` 和 `gate-2-design-review.md`。如果没有 review 文件，提示用户先执行 `/design:review`。
3. **读取门禁反馈**：读取 `gate-2-design-review.md` 中指出的不通过项和修改建议。收集用户附加的额外修改意见。
4. **针对性修正设计文档 (`design.md`)**：
   - 检查核心章节是否完整：本次需求的改动服务范围、整体流程图或流程说明、架构方案（包含改动代码结构）、可能带来的性能瓶颈点分析、潜在的风险点及应对措施。补充缺失部分。
   - 修正与服务归属和通信协议不符的内容。如果是普通HTTP Service，移除 gRPC/proto 相关章节，严格对齐 `code-implementation-examples.md` 中的 HTTP Service 样板。
   - 修正 gRPC service 的设计，确保符合 protobuf 规范和 `code-implementation-examples.md` 的样板。
   - 修正数据访问层设计，确保符合 DAO 样板。
   - 修正表结构设计，确保符合 db设计规范（仅建议索引，避免多列索引，生成的DDL不能有除主键外其他索引）。
5. **同步修正任务拆解**：基于修正后的 `design.md`，重新生成或更新 `tasks/features.json` 进行任务拆解。
6. **引导用户下一步**：提示用户重新执行 `/design:review`。
