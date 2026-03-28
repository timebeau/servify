#!/usr/bin/env bash
# demo-setup.sh — 一键启动 Servify 演示环境
# 启动服务 → 等待健康检查 → 填充演示数据 → 打印访问信息
set -euo pipefail

BOLD='\033[1m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BOLD}╔══════════════════════════════════════════╗${NC}"
echo -e "${BOLD}║     Servify 一键演示环境启动脚本        ║${NC}"
echo -e "${BOLD}╚══════════════════════════════════════════╝${NC}"
echo ""

BASE_URL="${1:-http://localhost:8080}"

# ── Step 1: 检查依赖 ──────────────────────────────────────
echo -e "${BLUE}[1/4] 检查依赖...${NC}"
for cmd in go curl; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "ERROR: $cmd 未安装，请先安装。"
    exit 1
  fi
done
echo -e "  ${GREEN}✓${NC} 依赖检查通过"

# ── Step 2: 构建 ──────────────────────────────────────────
echo -e "${BLUE}[2/4] 构建服务...${NC}"
if ! go build -o ./bin/servify ./apps/server/cmd/server 2>/dev/null; then
  echo "ERROR: 构建失败，请检查代码。"
  exit 1
fi
echo -e "  ${GREEN}✓${NC} 服务构建成功"

# ── Step 3: 启动服务（后台） ───────────────────────────────
echo -e "${BLUE}[3/4] 启动服务...${NC}"

# 检查是否已有实例运行
if curl -sfS "$BASE_URL/health" &>/dev/null; then
  echo -e "  ${YELLOW}⚠${NC} 检测到服务已在运行 ($BASE_URL)"
else
  export SERVIFY_JWT_SECRET="${SERVIFY_JWT_SECRET:-demo-secret-change-in-production}"
  export SERVIFY_DB_DSN="${SERVIFY_DB_DSN:-servify.db}"
  export SERVIFY_PORT="${SERVIFY_PORT:-8080}"

  ./bin/servify &
  SERVER_PID=$!
  echo "  服务 PID: $SERVER_PID"

  # 等待服务启动
  echo -n "  等待服务就绪"
  READY=0
  for i in $(seq 1 30); do
    if curl -sfS "$BASE_URL/health" &>/dev/null; then
      READY=1
      break
    fi
    echo -n "."
    sleep 1
  done
  if [ "$READY" -eq 1 ]; then
    echo ""
    echo -e "  ${GREEN}✓${NC} 服务已启动 ($BASE_URL)"
  else
    echo ""
    echo "ERROR: 服务启动超时（30秒）。请检查日志。"
    kill "$SERVER_PID" 2>/dev/null || true
    exit 1
  fi
fi

# ── Step 4: 填充演示数据 ──────────────────────────────────
echo -e "${BLUE}[4/4] 填充演示数据...${NC}"
chmod +x scripts/seed-data.sh
if scripts/seed-data.sh "$BASE_URL"; then
  echo -e "  ${GREEN}✓${NC} 演示数据填充完成"
else
  echo -e "  ${YELLOW}⚠${NC} 数据填充部分失败（可能数据已存在）"
fi

# ── 完成 ──────────────────────────────────────────────────
echo ""
echo -e "${BOLD}════════════════════════════════════════════${NC}"
echo -e "${GREEN}${BOLD}  演示环境已就绪！${NC}"
echo -e "${BOLD}════════════════════════════════════════════${NC}"
echo ""
echo -e "  ${BOLD}管理后台:${NC}  $BASE_URL/admin/"
echo -e "  ${BOLD}演示网站:${NC}  $BASE_URL/demo/ (需要 admin 构建后访问)"
echo -e "  ${BOLD}健康检查:${NC}  $BASE_URL/health"
echo ""
echo -e "  ${BOLD}测试账号:${NC}"
echo -e "    管理员:  admin / admin123"
echo -e "    客服:    agent1 / agent123"
echo -e "    客户:    customer1 / cust123"
echo ""
echo -e "  ${YELLOW}按 Ctrl+C 停止服务${NC}"
echo ""

# 保持前台运行
wait 2>/dev/null || true
