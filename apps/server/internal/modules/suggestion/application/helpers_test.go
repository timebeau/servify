package application_test

import (
	"testing"

	suggestionapp "servify/apps/server/internal/modules/suggestion/application"
)

func TestExtractTokens(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected []string
	}{
		{
			name:     "empty query",
			query:    "",
			expected: nil,
		},
		{
			name:     "whitespace query",
			query:    "   ",
			expected: nil,
		},
		{
			name:     "simple words",
			query:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			name:     "chinese characters - single char tokens",
			query:    "测试查询",
			expected: []string{"测", "试", "查", "询"},
		},
		{
			name:     "mixed",
			query:    "test 测试 api",
			expected: []string{"test", "测", "试", "api"},
		},
		{
			name:     "duplicates",
			query:    "test test api",
			expected: []string{"test", "api"},
		},
		{
			name:     "long token filtered",
			query:    "short " + string(make([]byte, 35)) + " another",
			expected: []string{"short", "another"},
		},
		{
			name:     "max tokens",
			query:    "one two three four five six seven eight nine ten eleven twelve thirteen fourteen fifteen sixteen seventeen eighteen nineteen twenty twentyone",
			expected: []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen", "twenty"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := suggestionapp.ExtractTokens(tt.query)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("ExtractTokens(%q) = %v, want nil", tt.query, result)
				}
				return
			}
			if len(result) != len(tt.expected) {
				t.Errorf("ExtractTokens(%q) length = %d, want %d", tt.query, len(result), len(tt.expected))
				return
			}
			for i, exp := range tt.expected {
				if result[i] != exp {
					t.Errorf("ExtractTokens(%q)[%d] = %q, want %q", tt.query, i, result[i], exp)
				}
			}
		})
	}
}

func TestClassifyIntent(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		expectedLabel string
		minConfidence float64
	}{
		{
			name:          "empty query",
			query:         "",
			expectedLabel: "general",
			minConfidence: 0.2,
		},
		{
			name:          "complaint chinese",
			query:         "我要投诉",
			expectedLabel: "complaint",
			minConfidence: 0.6,
		},
		{
			name:          "complaint english",
			query:         "I want to complaint",
			expectedLabel: "complaint",
			minConfidence: 0.6,
		},
		{
			name:          "billing chinese",
			query:         "怎么开发票",
			expectedLabel: "billing",
			minConfidence: 0.6,
		},
		{
			name:          "billing english",
			query:         "Need invoice",
			expectedLabel: "billing",
			minConfidence: 0.6,
		},
		{
			name:          "technical chinese",
			query:         "系统报错",
			expectedLabel: "technical",
			minConfidence: 0.6,
		},
		{
			name:          "technical english",
			query:         "System error",
			expectedLabel: "technical",
			minConfidence: 0.6,
		},
		{
			name:          "general query",
			query:         "Hello world",
			expectedLabel: "general",
			minConfidence: 0.3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := suggestionapp.ClassifyIntent(tt.query)
			if result.Label != tt.expectedLabel {
				t.Errorf("ClassifyIntent(%q) label = %q, want %q", tt.query, result.Label, tt.expectedLabel)
			}
			if result.Confidence < tt.minConfidence {
				t.Errorf("ClassifyIntent(%q) confidence = %f, want >= %f", tt.query, result.Confidence, tt.minConfidence)
			}
		})
	}
}

func TestScoreText(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		text     string
		expected float64
	}{
		{
			name:     "exact match",
			query:    "hello world",
			text:     "hello world",
			expected: 1.0,
		},
		{
			name:     "partial match",
			query:    "hello world",
			text:     "hello there",
			expected: 0.5,
		},
		{
			name:     "no match",
			query:    "hello",
			text:     "world",
			expected: 0.0,
		},
		{
			name:     "empty query",
			query:    "",
			text:     "hello",
			expected: 0.0,
		},
		{
			name:     "empty text",
			query:    "hello",
			text:     "",
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := suggestionapp.ScoreText(tt.query, tt.text)
			if result != tt.expected {
				t.Errorf("ScoreText(%q, %q) = %f, want %f", tt.query, tt.text, result, tt.expected)
			}
		})
	}
}

func TestBuildLikeWhereTokens(t *testing.T) {
	tests := []struct {
		name      string
		fields    []string
		tokens    []string
		maxTokens int
		wantEmpty bool
	}{
		{
			name:      "empty fields",
			fields:    []string{},
			tokens:    []string{"test"},
			maxTokens: 3,
			wantEmpty: true,
		},
		{
			name:      "empty tokens",
			fields:    []string{"title"},
			tokens:    []string{},
			maxTokens: 3,
			wantEmpty: true,
		},
		{
			name:      "zero max tokens",
			fields:    []string{"title"},
			tokens:    []string{"test"},
			maxTokens: 0,
			wantEmpty: true,
		},
		{
			name:      "valid",
			fields:    []string{"title", "content"},
			tokens:    []string{"test", "api"},
			maxTokens: 3,
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			where, args := suggestionapp.BuildLikeWhereTokens(tt.fields, tt.tokens, tt.maxTokens)
			isEmpty := where == "" || len(args) == 0
			if isEmpty != tt.wantEmpty {
				t.Errorf("BuildLikeWhereTokens() isEmpty = %v, want %v", isEmpty, tt.wantEmpty)
			}
		})
	}
}
