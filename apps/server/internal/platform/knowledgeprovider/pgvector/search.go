package pgvector

import (
	"fmt"
	"math"
	"strings"
)

// CosineDistance 计算两个向量之间的余弦距离
// 余弦距离 = 1 - 余弦相似度
// 余弦相似度 = (a·b) / (||a|| * ||b||)
func CosineDistance(a, b []float32) float32 {
	if len(a) != len(b) {
		panic("vectors must have the same dimension")
	}
	if len(a) == 0 {
		return 0
	}

	var dotProduct float32
	var normA float32
	var normB float32

	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	normA = sqrt32(normA)
	normB = sqrt32(normB)

	if normA == 0 || normB == 0 {
		return 0
	}

	cosineSim := dotProduct / (normA * normB)
	return 1 - cosineSim
}

// EuclideanDistance 计算欧几里得距离（L2 距离）
// d(a,b) = sqrt(sum((a[i] - b[i])^2))
func EuclideanDistance(a, b []float32) float32 {
	if len(a) != len(b) {
		panic("vectors must have the same dimension")
	}
	if len(a) == 0 {
		return 0
	}

	var sum float32
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}

	return sqrt32(sum)
}

// sqrt32 计算 float32 的平方根
// 使用 math.Sqrt 然后转换回 float32 以保持精度一致
func sqrt32(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

// BuildCosineSearchSQL 构建 pgvector 余弦相似度搜索 SQL
// 返回的 SQL 使用 <=> 操作符进行余弦距离计算
// 结果按距离升序排列，距离越小表示越相似
func BuildCosineSearchSQL(tableName string, dimension int) string {
	return fmt.Sprintf(`
SELECT id, content, metadata, embedding, distance
FROM %s,
     (SELECT $1::vector(%d) AS query) q
ORDER BY embedding <=> q.query
LIMIT $2;
`, tableName, dimension)
}

// BuildEuclideanSearchSQL 构建 pgvector 欧几里得距离搜索 SQL
// 返回的 SQL 使用 <-> 操作符进行 L2 距离计算
// 结果按距离升序排列，距离越小表示越相似
func BuildEuclideanSearchSQL(tableName string) string {
	return fmt.Sprintf(`
SELECT id, content, metadata, embedding, distance
FROM %s,
     (SELECT $1 AS query) q
ORDER BY embedding <-> q.query
LIMIT $2;
`, tableName)
}

// FormatVector 将向量转换为 pgvector 格式字符串
// pgvector 格式: '[0.1,0.2,0.3,...]'
func FormatVector(v []float32) string {
	if len(v) == 0 {
		return "[]"
	}

	var sb strings.Builder
	sb.Grow(len(v)*8 + 2) // 预分配足够空间

	sb.WriteRune('[')
	for i, val := range v {
		if i > 0 {
			sb.WriteRune(',')
		}
		// 使用 fmt 格式化保证精度
		sb.WriteString(fmt.Sprintf("%g", val))
	}
	sb.WriteRune(']')

	return sb.String()
}

// NormalizeVector L2 归一化向量
// 使向量的 L2 范数为 1
// normalized[i] = v[i] / sqrt(sum(v[j]^2))
func NormalizeVector(v []float32) []float32 {
	if len(v) == 0 {
		return v
	}

	var sum float32
	for _, val := range v {
		sum += val * val
	}

	norm := sqrt32(sum)
	if norm == 0 {
		// 零向量无法归一化，返回原向量
		return v
	}

	result := make([]float32, len(v))
	for i, val := range v {
		result[i] = val / norm
	}

	return result
}

// DotProduct 计算两个向量的点积（内积）
func DotProduct(a, b []float32) float32 {
	if len(a) != len(b) {
		panic("vectors must have the same dimension")
	}
	if len(a) == 0 {
		return 0
	}

	var result float32
	for i := range a {
		result += a[i] * b[i]
	}

	return result
}

// L2Norm 计算向量的 L2 范数
func L2Norm(v []float32) float32 {
	if len(v) == 0 {
		return 0
	}

	var sum float32
	for _, val := range v {
		sum += val * val
	}

	return sqrt32(sum)
}

// CosineSimilarity 计算两个向量的余弦相似度
// 返回值范围 [-1, 1]，1 表示完全相同方向，-1 表示完全相反方向
func CosineSimilarity(a, b []float32) float32 {
	return 1 - CosineDistance(a, b)
}
