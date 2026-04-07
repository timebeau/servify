#!/usr/bin/env bash
# seed-data.sh — 填充演示数据
# 前置条件：Servify 服务已启动（localhost:8080）
set -euo pipefail

BASE_URL="${1:-http://localhost:8080}"
echo "==> Seeding demo data to $BASE_URL"

extract_json_field() {
  local key="$1"
  python3 -c 'import json,sys
key = sys.argv[1]
try:
    data = json.load(sys.stdin)
except Exception:
    sys.exit(1)

def walk(value):
    if isinstance(value, dict):
        if key in value:
            result = value[key]
            if result is None:
                return None
            print(result)
            raise SystemExit(0)
        for child in value.values():
            walk(child)
    elif isinstance(value, list):
        for child in value:
            walk(child)

walk(data)
sys.exit(1)
' "$key"
}

# 注册 admin 用户（首个用户自动成为 admin）
echo "  Creating admin user..."
ADMIN_RESP=$(curl -sfS "$BASE_URL/api/v1/auth/register" \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","email":"admin@servify.io","password":"admin123","name":"管理员","role":"admin"}') || {
  # 如果已存在，尝试登录
  ADMIN_RESP=$(curl -sfS "$BASE_URL/api/v1/auth/login" \
    -H 'Content-Type: application/json' \
    -d '{"username":"admin","password":"admin123"}')
}
TOKEN=$(echo "$ADMIN_RESP" | grep -o '"token":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "ERROR: Failed to get admin token"
  exit 1
fi
echo "  Admin token obtained."

AUTH="Authorization: Bearer $TOKEN"

# 注册客服
echo "  Creating agent users..."
AGENT_USER_IDS=()
for i in 1 2 3; do
  AGENT_RESP=$(curl -sfS "$BASE_URL/api/v1/auth/register" \
    -H 'Content-Type: application/json' \
    -H "$AUTH" \
    -d "{\"username\":\"agent$i\",\"email\":\"agent$i@servify.io\",\"password\":\"agent123\",\"name\":\"客服$i\",\"role\":\"agent\"}" 2>/dev/null || true)

  if [ -z "$AGENT_RESP" ]; then
    AGENT_RESP=$(curl -sfS "$BASE_URL/api/v1/auth/login" \
      -H 'Content-Type: application/json' \
      -d "{\"username\":\"agent$i\",\"password\":\"agent123\"}")
  fi

  AGENT_USER_ID=$(echo "$AGENT_RESP" | extract_json_field id || true)
  if [ -n "$AGENT_USER_ID" ]; then
    AGENT_USER_IDS+=("$AGENT_USER_ID")
    curl -sfS "$BASE_URL/api/agents" \
      -H 'Content-Type: application/json' \
      -H "$AUTH" \
      -d "{\"user_id\":$AGENT_USER_ID,\"department\":\"客服\",\"skills\":\"工单,支持\",\"max_concurrent\":5}" > /dev/null 2>&1 || true
    curl -sfS -X POST "$BASE_URL/api/agents/$AGENT_USER_ID/online" \
      -H 'Content-Type: application/json' \
      -H "$AUTH" > /dev/null 2>&1 || true
  fi
done

# 注册客户
echo "  Creating customer users..."
CUSTOMER_USER_IDS=()
for i in 1 2 3 4 5; do
  CUSTOMER_RESP=$(curl -sfS "$BASE_URL/api/v1/auth/register" \
    -H 'Content-Type: application/json' \
    -d "{\"username\":\"customer$i\",\"email\":\"customer$i@example.com\",\"password\":\"cust123\",\"name\":\"客户$i\",\"role\":\"customer\"}" 2>/dev/null || true)

  if [ -z "$CUSTOMER_RESP" ]; then
    CUSTOMER_RESP=$(curl -sfS "$BASE_URL/api/v1/auth/login" \
      -H 'Content-Type: application/json' \
      -d "{\"username\":\"customer$i\",\"password\":\"cust123\"}")
  fi

  CUSTOMER_USER_ID=$(echo "$CUSTOMER_RESP" | extract_json_field id || true)
  if [ -n "$CUSTOMER_USER_ID" ]; then
    CUSTOMER_USER_IDS+=("$CUSTOMER_USER_ID")
  fi
done

# 创建工单
echo "  Creating demo tickets..."
if [ "${#CUSTOMER_USER_IDS[@]}" -lt 5 ]; then
  echo "ERROR: Failed to resolve customer user ids"
  exit 1
fi

TICKET_DATA=(
  '{"title":"无法登录账号","description":"输入正确密码后提示密码错误","priority":"high","category":"账号问题"}'
  '{"title":"退款申请","description":"订单 #12345 商品损坏，申请退款","priority":"urgent","category":"退款"}'
  '{"title":"产品使用咨询","description":"请问如何导出报表？","priority":"medium","category":"使用指导"}'
  '{"title":"功能建议","description":"建议增加批量导入功能","priority":"low","category":"建议"}'
  '{"title":"支付失败","description":"支付时提示网络超时","priority":"high","category":"支付"}'
  '{"title":"API 调用报错","description":"调用 /api/v1/tickets 返回 500","priority":"urgent","category":"技术支持"}'
  '{"title":"会员续费问题","description":"自动续费未生效","priority":"medium","category":"会员"}'
  '{"title":"配送延迟","description":"订单已超过预计到达时间 3 天","priority":"high","category":"物流"}'
  '{"title":"隐私设置咨询","description":"如何关闭个性化推荐","priority":"low","category":"隐私"}'
  '{"title":"发票申请","description":"需要开具增值税专用发票","priority":"medium","category":"发票"}'
)

for idx in "${!TICKET_DATA[@]}"; do
  DATA="${TICKET_DATA[$idx]}"
  CUSTOMER_ID="${CUSTOMER_USER_IDS[$((idx % 5))]}"
  PAYLOAD=$(python3 -c 'import json,sys; data=json.loads(sys.argv[1]); data["customer_id"]=int(sys.argv[2]); print(json.dumps(data, ensure_ascii=False))' "$DATA" "$CUSTOMER_ID")
  curl -sfS "$BASE_URL/api/tickets" \
    -H 'Content-Type: application/json' \
    -H "$AUTH" \
    -d "$PAYLOAD" > /dev/null 2>&1 || true
done

# 创建知识库文档
echo "  Creating knowledge docs..."
for data in \
  '{"title":"产品安装指南","content":"## 安装步骤\n\n1. 下载安装包\n2. 运行安装程序\n3. 配置数据库连接\n4. 启动服务","category":"技术支持"}' \
  '{"title":"常见问题解答","content":"## FAQ\n\n### Q: 如何重置密码？\nA: 点击登录页的「忘记密码」链接。\n\n### Q: 如何升级套餐？\nA: 进入设置 > 订阅管理 > 选择新套餐。","category":"常见问题"}' \
  '{"title":"API 使用文档","content":"## REST API\n\nBase URL: `/api/v1`\n\n认证方式: Bearer Token (JWT)\n\n主要端点:\n- POST /auth/login\n- GET /tickets\n- GET /customers","category":"开发文档"}' \
  '{"title":"退款政策","content":"## 退款政策\n\n- 购买后 7 天内可申请全额退款\n- 超过 7 天按比例退款\n- 虚拟商品不支持退款","category":"政策"}' \
  '{"title":"隐私政策","content":"## 隐私政策\n\n我们重视您的隐私。收集的数据仅用于改善服务质量。\n\n您可以在设置中管理隐私选项。","category":"政策"}'; do
  curl -sfS "$BASE_URL/api/knowledge-docs" \
    -H 'Content-Type: application/json' \
    -H "$AUTH" \
    -d "$data" > /dev/null 2>&1 || true
done

echo ""
echo "✅ Seed data complete!"
echo "   Admin:  admin / admin123"
echo "   Agents: agent1-3 / agent123"
echo "   Customers: customer1-5 / cust123"
echo "   Demo tickets: 10"
echo "   Knowledge docs: 5"
