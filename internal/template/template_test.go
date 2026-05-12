package template

import (
	"testing"

	"github.com/yuan-shuo/gometrics/internal/config"
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

func TestToCamel(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"snake_case", "snakeCase"},
		{"multiple_words_here", "multipleWordsHere"},
		{"", ""},
		{"a", "a"},
		{"already_Pascal", "alreadyPascal"},
		{"with_123_numbers", "with123Numbers"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToCamel(tt.input)
			if result != tt.expected {
				t.Errorf("ToCamel(%q) = %q, want %q", tt.input, result, tt.expected)
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

func TestJoin(t *testing.T) {
	tests := []struct {
		strs     []string
		sep      string
		expected string
	}{
		{[]string{"a", "b", "c"}, ", ", "a, b, c"},
		{[]string{"single"}, ", ", "single"},
		{[]string{}, ", ", ""},
		{nil, ", ", ""},
		{[]string{"a", "b"}, " | ", "a | b"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := Join(tt.strs, tt.sep)
			if result != tt.expected {
				t.Errorf("Join(%v, %q) = %q, want %q", tt.strs, tt.sep, result, tt.expected)
			}
		})
	}
}

func TestIsEnum(t *testing.T) {
	tests := []struct {
		vals     []string
		expected bool
	}{
		{[]string{"a", "b", "c"}, true},
		{[]string{"app", "web"}, true},
		{[]string{"*"}, false},
		{[]string{"a", "*"}, false},
		{[]string{}, false},
		{nil, false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := IsEnum(tt.vals)
			if result != tt.expected {
				t.Errorf("IsEnum(%v) = %v, want %v", tt.vals, result, tt.expected)
			}
		})
	}
}

func TestEnumVars(t *testing.T) {
	tests := []struct {
		metricName string
		labelName  string
		vals       []string
		expected   string
	}{
		{"requests_total", "source", []string{"app", "web"}, "RequestsTotalSourceApp, RequestsTotalSourceWeb"},
		{"active_connections", "pool", []string{"default"}, "ActiveConnectionsPoolDefault"},
		{"request_duration_ms", "method", []string{"GET", "POST", "PUT"}, "RequestDurationMsMethodGET, RequestDurationMsMethodPOST, RequestDurationMsMethodPUT"},
	}

	for _, tt := range tests {
		t.Run(tt.metricName+"_"+tt.labelName, func(t *testing.T) {
			result := EnumVars(tt.metricName, tt.labelName, tt.vals)
			if result != tt.expected {
				t.Errorf("EnumVars(%q, %q, %v) = %q, want %q", tt.metricName, tt.labelName, tt.vals, result, tt.expected)
			}
		})
	}
}

func TestLabelParams(t *testing.T) {
	tests := []struct {
		metricName string
		labels     []config.Label
		expected   string
	}{
		{
			"requests_total",
			[]config.Label{
				{Name: "method", Vals: []string{"GET", "POST"}},
				{Name: "path", Vals: []string{"*"}},
			},
			"method RequestsTotalMethod, path string",
		},
		{
			"active_connections",
			[]config.Label{
				{Name: "pool", Vals: []string{"default", "cache"}},
			},
			"pool ActiveConnectionsPool",
		},
		{
			"simple_counter",
			[]config.Label{},
			"",
		},
		{
			"no_labels",
			nil,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.metricName, func(t *testing.T) {
			result := LabelParams(tt.metricName, tt.labels)
			if result != tt.expected {
				t.Errorf("LabelParams(%q, %v) = %q, want %q", tt.metricName, tt.labels, result, tt.expected)
			}
		})
	}
}

func TestLabelArg(t *testing.T) {
	tests := []struct {
		label      config.Label
		metricName string
		expected   string
	}{
		{
			config.Label{Name: "source", Vals: []string{"app", "web"}},
			"requests_total",
			"source.requestsTotalsourceValue()",
		},
		{
			config.Label{Name: "path", Vals: []string{"*"}},
			"requests_total",
			"path",
		},
		{
			config.Label{Name: "pool", Vals: []string{"default"}},
			"active_connections",
			"pool.activeConnectionspoolValue()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.label.Name, func(t *testing.T) {
			result := LabelArg(tt.label, tt.metricName)
			if result != tt.expected {
				t.Errorf("LabelArg(%v, %q) = %q, want %q", tt.label, tt.metricName, result, tt.expected)
			}
		})
	}
}

func TestFuncMap(t *testing.T) {
	fm := FuncMap()

	// 验证所有预期的函数都存在
	expectedFuncs := []string{"toPascal", "camelCase", "hasMethod", "isEnum", "join", "enumVars", "labelParams", "labelArg"}
	for _, name := range expectedFuncs {
		if _, ok := fm[name]; !ok {
			t.Errorf("FuncMap() missing function %q", name)
		}
	}
}
