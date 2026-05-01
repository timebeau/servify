#!/bin/bash
# scripts/test-knowledge-acceptance.sh
# 知识库验收脚本

set -e

SERVIFY_URL="${SERVIFY_URL:-http://localhost:8080}"
EMBEDDING_PROVIDER="${EMBEDDING_PROVIDER:-openai}"
EVIDENCE_DIR="${EVIDENCE_DIR:-./scripts/test-results/knowledge-acceptance}"

echo "=== Servify Knowledge Base Acceptance Test ==="
echo "Servify URL: $SERVIFY_URL"
echo "Embedding Provider: $EMBEDDING_PROVIDER"
echo "Evidence Directory: $EVIDENCE_DIR"
echo ""

# 创建证据目录
mkdir -p "$EVIDENCE_DIR"

# 1. 健康检查
echo "1. Health Check..."
curl -s "$SERVIFY_URL/health" | tee "$EVIDENCE_DIR/health.json"
echo ""

# 2. 创建测试文档
echo "2. Creating test document..."
curl -s -X POST "$SERVIFY_URL/api/knowledge-docs" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TEST_TOKEN" \
  -d '{
    "title": "Test Document",
    "content": "This is a test document for knowledge base acceptance testing. It contains information about product features, installation steps, and troubleshooting guides.",
    "category": "test",
    "tags": ["acceptance", "test"],
    "is_public": true
  }' | tee "$EVIDENCE_DIR/create-doc.json"
echo ""

# 3. 列出文档
echo "3. Listing documents..."
curl -s "$SERVIFY_URL/api/knowledge-docs?page=1&page_size=10" \
  -H "Authorization: Bearer $TEST_TOKEN" | tee "$EVIDENCE_DIR/list-docs.json"
echo ""

# 4. 搜索测试
echo "4. Searching documents..."
curl -s -X POST "$SERVIFY_URL/api/v1/ai/query" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "product features and installation",
    "session_id": "acceptance-test"
  }' | tee "$EVIDENCE_DIR/search-result.json"
echo ""

# 5. 生成 manifest
echo "5. Generating manifest..."
cat > "$EVIDENCE_DIR/manifest.json" << EOF
{
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "servify_url": "$SERVIFY_URL",
  "embedding_provider": "$EMBEDDING_PROVIDER",
  "tests": [
    {"name": "health", "file": "health.json"},
    {"name": "create_doc", "file": "create-doc.json"},
    {"name": "list_docs", "file": "list-docs.json"},
    {"name": "search", "file": "search-result.json"}
  ]
}
EOF

echo "=== Acceptance Test Complete ==="
echo "Evidence saved to: $EVIDENCE_DIR"
