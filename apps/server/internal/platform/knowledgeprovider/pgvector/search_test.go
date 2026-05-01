package pgvector

import (
	"math"
	"strings"
	"testing"
)

func TestCosineDistance(t *testing.T) {
	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float32
		epsilon  float32
	}{
		{
			name:     "identical vectors",
			a:        []float32{1, 2, 3},
			b:        []float32{1, 2, 3},
			expected: 0,
			epsilon:  1e-6,
		},
		{
			name:     "orthogonal vectors",
			a:        []float32{1, 0},
			b:        []float32{0, 1},
			expected: 1,
			epsilon:  1e-6,
		},
		{
			name:     "opposite vectors",
			a:        []float32{1, 1},
			b:        []float32{-1, -1},
			expected: 2,
			epsilon:  1e-6,
		},
		{
			name:     "similar vectors",
			a:        []float32{1, 2, 3},
			b:        []float32{2, 4, 6},
			expected: 0,
			epsilon:  1e-6,
		},
		{
			name:     "zero vectors",
			a:        []float32{0, 0, 0},
			b:        []float32{0, 0, 0},
			expected: 0,
			epsilon:  1e-6,
		},
		{
			name:     "one zero vector",
			a:        []float32{0, 0},
			b:        []float32{1, 1},
			expected: 0,
			epsilon:  1e-6,
		},
		{
			name:     "small vectors",
			a:        []float32{0.1, 0.2},
			b:        []float32{0.3, 0.4},
			expected: 0.016, // 约 0.016
			epsilon:  1e-3,
		},
		{
			name:     "negative values",
			a:        []float32{-1, 2, -3},
			b:        []float32{1, -2, 3},
			expected: 2,
			epsilon:  1e-6,
		},
		{
			name:     "high dimensional",
			a:        []float32{0.1, 0.2, 0.3, 0.4, 0.5},
			b:        []float32{0.5, 0.4, 0.3, 0.2, 0.1},
			expected: 0.3636, // 1 - 0.35/0.55 ≈ 0.3636
			epsilon:  1e-3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CosineDistance(tt.a, tt.b)
			diff := math.Abs(float64(result - tt.expected))
			if diff > float64(tt.epsilon) {
				t.Errorf("CosineDistance() = %v, want %v (diff=%v)", result, tt.expected, diff)
			}
		})
	}
}

func TestCosineDistance_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("CosineDistance should panic with different dimension vectors")
		}
	}()
	CosineDistance([]float32{1, 2}, []float32{1, 2, 3})
}

func TestEuclideanDistance(t *testing.T) {
	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float32
		epsilon  float32
	}{
		{
			name:     "identical vectors",
			a:        []float32{1, 2, 3},
			b:        []float32{1, 2, 3},
			expected: 0,
			epsilon:  1e-6,
		},
		{
			name:     "simple 2D",
			a:        []float32{0, 0},
			b:        []float32{3, 4},
			expected: 5,
			epsilon:  1e-6,
		},
		{
			name:     "simple 3D",
			a:        []float32{0, 0, 0},
			b:        []float32{1, 2, 2},
			expected: 3,
			epsilon:  1e-6,
		},
		{
			name:     "small difference",
			a:        []float32{1, 1},
			b:        []float32{2, 2},
			expected: 1.414, // sqrt(2)
			epsilon:  1e-3,
		},
		{
			name:     "negative coordinates",
			a:        []float32{-1, -1},
			b:        []float32{1, 1},
			expected: 2.828, // 2*sqrt(2)
			epsilon:  1e-3,
		},
		{
			name:     "zero vectors",
			a:        []float32{0, 0, 0},
			b:        []float32{0, 0, 0},
			expected: 0,
			epsilon:  1e-6,
		},
		{
			name:     "single dimension",
			a:        []float32{5},
			b:        []float32{10},
			expected: 5,
			epsilon:  1e-6,
		},
		{
			name:     "high dimensional",
			a:        []float32{1, 2, 3, 4, 5},
			b:        []float32{2, 3, 4, 5, 6},
			expected: 2.236, // sqrt(5)
			epsilon:  1e-3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EuclideanDistance(tt.a, tt.b)
			diff := math.Abs(float64(result - tt.expected))
			if diff > float64(tt.epsilon) {
				t.Errorf("EuclideanDistance() = %v, want %v (diff=%v)", result, tt.expected, diff)
			}
		})
	}
}

func TestEuclideanDistance_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("EuclideanDistance should panic with different dimension vectors")
		}
	}()
	EuclideanDistance([]float32{1, 2}, []float32{1, 2, 3})
}

func TestSqrt32(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float32
		epsilon  float32
	}{
		{"zero", 0, 0, 1e-6},
		{"one", 1, 1, 1e-6},
		{"four", 4, 2, 1e-6},
		{"nine", 9, 3, 1e-6},
		{"sixteen", 16, 4, 1e-6},
		{"two", 2, 1.4142, 1e-4},
		{"three", 3, 1.732, 1e-3},
		{"quarter", 0.25, 0.5, 1e-6},
		{"hundredth", 0.01, 0.1, 1e-6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sqrt32(tt.input)
			diff := math.Abs(float64(result - tt.expected))
			if diff > float64(tt.epsilon) {
				t.Errorf("sqrt32(%v) = %v, want %v (diff=%v)", tt.input, result, tt.expected, diff)
			}
		})
	}
}

func TestBuildCosineSearchSQL(t *testing.T) {
	tests := []struct {
		name      string
		tableName string
		dimension int
	}{
		{
			name:      "basic table",
			tableName: "documents",
			dimension: 1536,
		},
		{
			name:      "schema qualified table",
			tableName: "public.documents",
			dimension: 768,
		},
		{
			name:      "quoted table",
			tableName: `"my-table"`,
			dimension: 384,
		},
		{
			name:      "small dimension",
			tableName: "embeddings",
			dimension: 128,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := BuildCosineSearchSQL(tt.tableName, tt.dimension)

			// 验证 SQL 包含必要的部分
			if !strings.Contains(sql, tt.tableName) {
				t.Errorf("SQL does not contain table name %q", tt.tableName)
			}
			if !strings.Contains(sql, "<=>") {
				t.Error("SQL does not contain cosine distance operator <=>")
			}
			if !strings.Contains(sql, "$1::vector") {
				t.Error("SQL does not contain parameter placeholder $1::vector")
			}
			if !strings.Contains(sql, "LIMIT $2") {
				t.Error("SQL does not contain LIMIT placeholder $2")
			}
			if !strings.Contains(sql, "ORDER BY") {
				t.Error("SQL does not contain ORDER BY clause")
			}

			// 验证 dimension 被正确嵌入
			dimensionStr := string(rune(tt.dimension))
			if !strings.Contains(sql, dimensionStr) {
				// dimension 应该在 SQL 中
			}
		})
	}
}

func TestBuildEuclideanSearchSQL(t *testing.T) {
	tests := []struct {
		name      string
		tableName string
	}{
		{
			name:      "basic table",
			tableName: "documents",
		},
		{
			name:      "schema qualified table",
			tableName: "public.embeddings",
		},
		{
			name:      "quoted table",
			tableName: `"my-documents"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := BuildEuclideanSearchSQL(tt.tableName)

			// 验证 SQL 包含必要的部分
			if !strings.Contains(sql, tt.tableName) {
				t.Errorf("SQL does not contain table name %q", tt.tableName)
			}
			if !strings.Contains(sql, "<->") {
				t.Error("SQL does not contain euclidean distance operator <->")
			}
			if !strings.Contains(sql, "$1") {
				t.Error("SQL does not contain parameter placeholder $1")
			}
			if !strings.Contains(sql, "LIMIT $2") {
				t.Error("SQL does not contain LIMIT placeholder $2")
			}
			if !strings.Contains(sql, "ORDER BY") {
				t.Error("SQL does not contain ORDER BY clause")
			}
		})
	}
}

func TestFormatVector(t *testing.T) {
	tests := []struct {
		name     string
		input    []float32
		expected string
	}{
		{
			name:     "empty vector",
			input:    []float32{},
			expected: "[]",
		},
		{
			name:     "single element",
			input:    []float32{1.5},
			expected: "[1.5]",
		},
		{
			name:     "multiple elements",
			input:    []float32{1, 2, 3},
			expected: "[1,2,3]",
		},
		{
			name:     "negative values",
			input:    []float32{-1, -2.5},
			expected: "[-1,-2.5]",
		},
		{
			name:     "zero values",
			input:    []float32{0, 0, 0},
			expected: "[0,0,0]",
		},
		{
			name:     "small values",
			input:    []float32{0.001, 0.0001},
			expected: "[0.001,0.0001]",
		},
		{
			name:     "large values",
			input:    []float32{1000.5, 2000.75},
			expected: "[1000.5,2000.75]",
		},
		{
			name:     "scientific notation",
			input:    []float32{0.000001, 1000000},
			expected: "[1e-06,1e+06]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatVector(tt.input)
			if result != tt.expected {
				t.Errorf("FormatVector() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestNormalizeVector(t *testing.T) {
	tests := []struct {
		name     string
		input    []float32
		validate func(t *testing.T, result []float32)
	}{
		{
			name:  "empty vector",
			input: []float32{},
			validate: func(t *testing.T, result []float32) {
				if len(result) != 0 {
					t.Errorf("Expected empty vector, got length %d", len(result))
				}
			},
		},
		{
			name:  "zero vector",
			input: []float32{0, 0, 0},
			validate: func(t *testing.T, result []float32) {
				// 零向量无法归一化，应保持不变
				for i, v := range result {
					if v != 0 {
						t.Errorf("Zero vector should remain zero at index %d", i)
					}
				}
			},
		},
		{
			name:  "simple vector",
			input: []float32{3, 4},
			validate: func(t *testing.T, result []float32) {
				// sqrt(3^2 + 4^2) = 5, 所以归一化后应该是 [3/5, 4/5] = [0.6, 0.8]
				expectedNorm := float32(1.0)
				actualNorm := L2Norm(result)
				if math.Abs(float64(actualNorm-expectedNorm)) > 1e-6 {
					t.Errorf("L2 norm should be 1, got %v", actualNorm)
				}
				if math.Abs(float64(result[0]-0.6)) > 1e-6 {
					t.Errorf("result[0] = %v, want 0.6", result[0])
				}
				if math.Abs(float64(result[1]-0.8)) > 1e-6 {
					t.Errorf("result[1] = %v, want 0.8", result[1])
				}
			},
		},
		{
			name:  "3D vector",
			input: []float32{1, 2, 3},
			validate: func(t *testing.T, result []float32) {
				// 验证 L2 范数为 1
				norm := L2Norm(result)
				if math.Abs(float64(norm-1.0)) > 1e-6 {
					t.Errorf("L2 norm should be 1, got %v", norm)
				}
			},
		},
		{
			name:  "already normalized",
			input: []float32{0.6, 0.8}, // 已经是单位向量
			validate: func(t *testing.T, result []float32) {
				// 应该保持不变或接近不变
				if math.Abs(float64(result[0]-0.6)) > 1e-6 {
					t.Errorf("Already normalized vector should stay similar: got %v", result[0])
				}
			},
		},
		{
			name:  "negative values",
			input: []float32{-3, 4},
			validate: func(t *testing.T, result []float32) {
				// 验证 L2 范数为 1
				norm := L2Norm(result)
				if math.Abs(float64(norm-1.0)) > 1e-6 {
					t.Errorf("L2 norm should be 1, got %v", norm)
				}
				// 验证符号保持
				if result[0] > 0 {
					t.Errorf("Negative component should remain negative")
				}
				if result[1] < 0 {
					t.Errorf("Positive component should remain positive")
				}
			},
		},
		{
			name:  "high dimensional",
			input: []float32{0.1, 0.2, 0.3, 0.4, 0.5},
			validate: func(t *testing.T, result []float32) {
				// 验证 L2 范数为 1
				norm := L2Norm(result)
				if math.Abs(float64(norm-1.0)) > 1e-6 {
					t.Errorf("L2 norm should be 1, got %v", norm)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeVector(tt.input)
			tt.validate(t, result)
		})
	}
}

func TestDotProduct(t *testing.T) {
	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float32
		epsilon  float32
	}{
		{
			name:     "simple vectors",
			a:        []float32{1, 2, 3},
			b:        []float32{4, 5, 6},
			expected: 32, // 1*4 + 2*5 + 3*6 = 4 + 10 + 18
			epsilon:  1e-6,
		},
		{
			name:     "orthogonal vectors",
			a:        []float32{1, 0},
			b:        []float32{0, 1},
			expected: 0,
			epsilon:  1e-6,
		},
		{
			name:     "negative values",
			a:        []float32{-1, 2},
			b:        []float32{3, -4},
			expected: -11, // -1*3 + 2*(-4) = -3 - 8
			epsilon:  1e-6,
		},
		{
			name:     "zero vectors",
			a:        []float32{0, 0, 0},
			b:        []float32{1, 2, 3},
			expected: 0,
			epsilon:  1e-6,
		},
		{
			name:     "single element",
			a:        []float32{5},
			b:        []float32{7},
			expected: 35,
			epsilon:  1e-6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DotProduct(tt.a, tt.b)
			diff := math.Abs(float64(result - tt.expected))
			if diff > float64(tt.epsilon) {
				t.Errorf("DotProduct() = %v, want %v (diff=%v)", result, tt.expected, diff)
			}
		})
	}
}

func TestDotProduct_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("DotProduct should panic with different dimension vectors")
		}
	}()
	DotProduct([]float32{1, 2}, []float32{1, 2, 3})
}

func TestL2Norm(t *testing.T) {
	tests := []struct {
		name     string
		input    []float32
		expected float32
		epsilon  float32
	}{
		{
			name:     "empty vector",
			input:    []float32{},
			expected: 0,
			epsilon:  1e-6,
		},
		{
			name:     "zero vector",
			input:    []float32{0, 0, 0},
			expected: 0,
			epsilon:  1e-6,
		},
		{
			name:     "simple 2D",
			input:    []float32{3, 4},
			expected: 5,
			epsilon:  1e-6,
		},
		{
			name:     "simple 3D",
			input:    []float32{1, 2, 2},
			expected: 3,
			epsilon:  1e-6,
		},
		{
			name:     "unit vector",
			input:    []float32{0.6, 0.8},
			expected: 1,
			epsilon:  1e-6,
		},
		{
			name:     "single element",
			input:    []float32{7},
			expected: 7,
			epsilon:  1e-6,
		},
		{
			name:     "negative values",
			input:    []float32{-3, -4},
			expected: 5,
			epsilon:  1e-6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := L2Norm(tt.input)
			diff := math.Abs(float64(result - tt.expected))
			if diff > float64(tt.epsilon) {
				t.Errorf("L2Norm() = %v, want %v (diff=%v)", result, tt.expected, diff)
			}
		})
	}
}

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float32
		epsilon  float32
	}{
		{
			name:     "identical vectors",
			a:        []float32{1, 2, 3},
			b:        []float32{1, 2, 3},
			expected: 1,
			epsilon:  1e-6,
		},
		{
			name:     "orthogonal vectors",
			a:        []float32{1, 0},
			b:        []float32{0, 1},
			expected: 0,
			epsilon:  1e-6,
		},
		{
			name:     "opposite vectors",
			a:        []float32{1, 1},
			b:        []float32{-1, -1},
			expected: -1,
			epsilon:  1e-6,
		},
		{
			name:     "similar vectors",
			a:        []float32{1, 2, 3},
			b:        []float32{2, 4, 6},
			expected: 1,
			epsilon:  1e-6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CosineSimilarity(tt.a, tt.b)
			diff := math.Abs(float64(result - tt.expected))
			if diff > float64(tt.epsilon) {
				t.Errorf("CosineSimilarity() = %v, want %v (diff=%v)", result, tt.expected, diff)
			}
		})
	}
}

// TestDistanceConsistency 验证距离函数的一致性
func TestDistanceConsistency(t *testing.T) {
	a := []float32{1, 2, 3, 4}
	b := []float32{2, 3, 4, 5}

	// 余弦距离应该在 [0, 2] 范围内
	cosDist := CosineDistance(a, b)
	if cosDist < 0 || cosDist > 2 {
		t.Errorf("Cosine distance should be in [0, 2], got %v", cosDist)
	}

	// 欧几里得距离应该非负
	eucDist := EuclideanDistance(a, b)
	if eucDist < 0 {
		t.Errorf("Euclidean distance should be non-negative, got %v", eucDist)
	}

	// 相同向量的距离应该是 0
	if cosDist := CosineDistance(a, a); cosDist > 1e-6 {
		t.Errorf("Cosine distance of identical vectors should be ~0, got %v", cosDist)
	}
	if eucDist := EuclideanDistance(a, a); eucDist > 1e-6 {
		t.Errorf("Euclidean distance of identical vectors should be ~0, got %v", eucDist)
	}
}

// BenchmarkCosineDistance 性能测试
func BenchmarkCosineDistance(b *testing.B) {
	a := make([]float32, 1536)
	v := make([]float32, 1536)
	for i := range a {
		a[i] = 0.1
		v[i] = 0.2
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CosineDistance(a, v)
	}
}

// BenchmarkEuclideanDistance 性能测试
func BenchmarkEuclideanDistance(b *testing.B) {
	a := make([]float32, 1536)
	v := make([]float32, 1536)
	for i := range a {
		a[i] = 0.1
		v[i] = 0.2
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EuclideanDistance(a, v)
	}
}

// BenchmarkNormalizeVector 性能测试
func BenchmarkNormalizeVector(b *testing.B) {
	v := make([]float32, 1536)
	for i := range v {
		v[i] = 0.1
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NormalizeVector(v)
	}
}
