package template

import (
	"testing"
)

func TestToPascal(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "Simple"},
		{"snake_case", "SnakeCase"},
		{"multiple_words_here", "MultipleWordsHere"},
		{"", ""},
		{"a", "A"},
		{"already_Pascal", "AlreadyPascal"},
		{"with_123_numbers", "With123Numbers"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToPascal(tt.input)
			if result != tt.expected {
				t.Errorf("ToPascal(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestHasMethod(t *testing.T) {
	tests := []struct {
		methods  []string
		target   string
		expected bool
	}{
		{[]string{"inc", "add"}, "inc", true},
		{[]string{"inc", "add"}, "add", true},
		{[]string{"inc", "add"}, "set", false},
		{[]string{}, "inc", false},
		{nil, "inc", false},
		{[]string{"set", "inc", "dec"}, "inc", true},
	}

	for _, tt := range tests {
		t.Run(tt.target, func(t *testing.T) {
			result := HasMethod(tt.methods, tt.target)
			if result != tt.expected {
				t.Errorf("HasMethod(%v, %q) = %v, want %v", tt.methods, tt.target, result, tt.expected)
			}
		})
	}
}

func TestJoinLabels(t *testing.T) {
	tests := []struct {
		labels   []string
		expected string
	}{
		{[]string{"a", "b", "c"}, "a, b, c"},
		{[]string{"single"}, "single"},
		{[]string{}, ""},
		{nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := JoinLabels(tt.labels)
			if result != tt.expected {
				t.Errorf("JoinLabels(%v) = %q, want %q", tt.labels, result, tt.expected)
			}
		})
	}
}

func TestLabelParams(t *testing.T) {
	tests := []struct {
		labels   []string
		expected string
	}{
		{[]string{"method", "path"}, "method string, path string"},
		{[]string{"single"}, "single string"},
		{[]string{}, ""},
		{nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := LabelParams(tt.labels)
			if result != tt.expected {
				t.Errorf("LabelParams(%v) = %q, want %q", tt.labels, result, tt.expected)
			}
		})
	}
}

func TestLabelArgs(t *testing.T) {
	tests := []struct {
		labels   []string
		expected string
	}{
		{[]string{"method", "path"}, "method, path"},
		{[]string{"single"}, "single"},
		{[]string{}, ""},
		{nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := LabelArgs(tt.labels)
			if result != tt.expected {
				t.Errorf("LabelArgs(%v) = %q, want %q", tt.labels, result, tt.expected)
			}
		})
	}
}

func TestFuncMap(t *testing.T) {
	fm := FuncMap()

	// 验证所有预期的函数都存在
	expectedFuncs := []string{"toPascal", "hasMethod", "joinLabels", "labelParams", "labelArgs"}
	for _, name := range expectedFuncs {
		if _, ok := fm[name]; !ok {
			t.Errorf("FuncMap() missing function %q", name)
		}
	}
}
