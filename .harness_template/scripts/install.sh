#!/usr/bin/env bash
# =============================================================================
# scripts/install.sh — Harness Engineering 初始化脚本（Claude Code 专用）
#
# 说明：
#   本项目使用 Claude Code 作为 AI 编程工具。
#   .claude/ 目录是 Claude Code 的原生配置目录，无需额外渲染。
#
#   本脚本负责：
#   1. 验证 .claude/ 目录结构完整性
#   2. 验证 .service-matrix/dependencies.yaml 格式
#   3. 初始化本地配置（.harness/local.yaml）
#
# 使用：
#   bash scripts/install.sh
#   bash scripts/install.sh --validate-only   # 只验证，不修改
# =============================================================================

set -euo pipefail

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

VALIDATE_ONLY=false
while [[ $# -gt 0 ]]; do
    case $1 in
        --validate-only) VALIDATE_ONLY=true; shift ;;
        -h|--help)
            echo "Usage: $0 [--validate-only]"
            exit 0
            ;;
        *) echo -e "${RED}Unknown option: $1${NC}"; exit 1 ;;
    esac
done

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
log_ok() { echo -e "${GREEN}  ✅${NC} $1"; }
log_fail() { echo -e "${RED}  ❌${NC} $1"; }

ERRORS=0

echo ""
echo "============================================"
echo "  Harness Engineering — 初始化验证"
echo "  (Claude Code 专用)"
echo "============================================"
echo ""

# ============================================================
# Step 1: 验证 Claude Code 目录结构
# ============================================================
log_info "Step 1: 验证 .claude/ 目录结构..."

check_file() {
    local file="$1"
    local desc="$2"
    if [[ -f "$PROJECT_ROOT/$file" ]]; then
        log_ok "$desc ($file)"
    else
        log_fail "$desc 缺失：$file"
        ERRORS=$((ERRORS + 1))
    fi
}

check_file "CLAUDE.md" "项目规范文件"
check_file ".claude/settings.json" "Claude 安全与权限配置"
check_file ".claude/commands/release-start.md" "版本新建命令"
check_file ".claude/commands/release-full.md" "全量发布梳理命令"
check_file ".claude/commands/release-incremental.md" "增量发布梳理命令"
check_file ".claude/commands/requirement-new.md" "需求新建命令"
check_file ".claude/commands/requirement-write.md" "需求撰写命令"
check_file ".claude/commands/requirement-start.md" "需求新建与撰写整合命令"
check_file ".claude/commands/agentic-code-review.md" "代码审查命令"
check_file ".claude/commands/knowledge-extract-experience.md" "经验提取命令"
check_file ".claude/commands/service-deps.md" "服务依赖命令"
check_file ".claude/skills/managing-requirement-lifecycle.md" "需求生命周期 Skill"
check_file ".claude/skills/code-review-report.md" "代码审查 Skill"
check_file ".claude/skills/self-refinement.md" "Self-Refinement Skill"
check_file ".claude/skills/service-dependency-analyzer.md" "服务依赖分析 Skill"
check_file ".claude/skills/traceability-gate-checker.md" "追溯链校验 Skill"
check_file ".claude/agents/Definition/requirement-quality-reviewer.md" "需求评审 Agent"
check_file ".claude/agents/DetailDesign/detail-design-quality-reviewer.md" "设计评审 Agent"
check_file ".claude/agents/Implementation/code-review-preparer.md" "代码审查准备 Agent"

# ============================================================
# Step 2: 验证核心工程制品
# ============================================================
echo ""
log_info "Step 2: 验证核心工程制品..."

check_file "context/harness-framework/main-process-numbering.md" "五阶段+四门禁（唯一真相源）"
check_file "context/team/INDEX.md" "团队级知识库"
check_file "context/team/git-convention.md" "Git 规范"
check_file "context/team/error-code.md" "错误码规范"
check_file "context/team/logging.md" "日志规范"
check_file ".service-matrix/dependencies.yaml" "服务矩阵"
check_file "releases/INDEX.md" "版本管理目录"
check_file "releases/template/RELEASE_NOTES.md" "版本说明模板"
check_file "releases/template/checklist.md" "发布检查清单模板"

# ============================================================
# Step 3: 初始化本地配置（非 validate-only 模式）
# ============================================================
if [[ "$VALIDATE_ONLY" == "false" ]]; then
    echo ""
    log_info "Step 3: 初始化本地配置..."

    # 创建 .harness/local.yaml（本地配置，不提交）
    HARNESS_LOCAL="$PROJECT_ROOT/.harness/local.yaml"
    if [[ ! -f "$HARNESS_LOCAL" ]]; then
        mkdir -p "$(dirname "$HARNESS_LOCAL")"
        cat > "$HARNESS_LOCAL" << EOF
# .harness/local.yaml — Harness Engineering 本地个人配置
#
# ⚠️ 此文件已加入 .gitignore，不会提交到 git 仓库
# ⚠️ 每个开发者在本机有自己的独立配置，互不影响
#
# 作用：
#   解决"占位符 → 真实路径"的本地绑定问题。
#
#   .service-matrix/dependencies.yaml 中使用占位符（如 {repo-root}）
#   而不是硬编码路径，原因是：
#   - 不同开发者的本地目录结构不同
#   - 不同机器（开发机/CI/生产机）的路径不同
#   - 路径不应该进入版本控制
#
#   本文件就是"将占位符替换为本机真实路径"的配置文件。
#
# 使用方式：
#   1. 运行 bash scripts/install.sh 时，会自动创建此文件模板
#   2. 根据本机实际情况填写下面的配置项
#   3. Claude Code 在执行命令时会读取此文件解析占位符

# -------------------------------------------------------------------
# 仓库根目录路径（必填）
# 说明：本仓库（monorepo-layout）在本机的绝对路径
# 用于替换 .service-matrix/dependencies.yaml 中的 {repo-root} 占位符
# 示例：
# repo_root: /Users/johndoe/workspace/api

# -------------------------------------------------------------------
# 当前活跃团队（可选，覆盖 .service-matrix 中的 default_team）
# 说明：当你同时参与多个团队的工作时，用于快速切换上下文
# 支持通过环境变量覆盖：export HARNESS_TEAM=another-team
# 示例：
# active_team: api

# -------------------------------------------------------------------
# 个人偏好设置（可选）
preferences:
  # 代码审查时，优先显示高严重性问题
  code_review_show_critical_first: true

  # Self-Refinement 时，自动推荐沉淀层级（无需每次确认）
  self_refinement_auto_suggest: true

  # 需求门禁检查时，显示详细的检查步骤
  gate_check_verbose: false
EOF
        log_ok "已创建本地配置：.harness/local.yaml"
    else
        log_ok "本地配置已存在：.harness/local.yaml"
    fi
fi

# ============================================================
# 结果汇总
# ============================================================
echo ""
echo "============================================"
if [[ $ERRORS -eq 0 ]]; then
    echo -e "${GREEN}  ✅ 验证通过！所有文件完整。${NC}"
    echo ""
    echo "下一步："
    echo "  1. 填写 .service-matrix/dependencies.yaml 中的服务信息"
    echo "  2. 更新 .harness/local.yaml 中的本地路径"
    echo "  3. 在 Claude Code 中输入 /requirement:start 开始第一个需求"
else
    echo -e "${RED}  ❌ 验证失败！发现 $ERRORS 个问题。${NC}"
    echo ""
    echo "请检查上述缺失文件，然后重新运行此脚本。"
    exit 1
fi
echo "============================================"
echo ""
