package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yuan-shuo/gometrics/internal/config"
)

func TestNew(t *testing.T) {
	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if gen == nil {
		t.Fatal("New() returned nil")
	}
	if gen.tmpl == nil {
		t.Fatal("New() returned generator with nil template")
	}
}

func TestGenerate(t *testing.T) {
	// 创建临时输出目录
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "testmetrics")

	// 创建测试配置 - 使用新的扁平结构
	cfg := &config.MetricConfig{
		Service: "test-service",
		Metrics: []config.Metric{
			{
				Name:    "requests_total",
				Help:    "Total requests",
				Type:    "counter",
				Labels:  []config.Label{{Name: "method", Vals: []string{"GET", "POST"}}},
				Methods: []string{"inc"},
			},
			{
				Name:    "active_connections",
				Help:    "Active connections",
				Type:    "gauge",
				Labels:  []config.Label{{Name: "pool", Vals: []string{"default", "cache"}}},
				Methods: []string{"set", "inc", "dec"},
			},
			{
				Name:    "request_duration_ms",
				Help:    "Request duration",
				Type:    "histogram",
				Labels:  []config.Label{{Name: "path", Vals: []string{"*"}}},
				Methods: []string{"observe"},
				Buckets: []float64{10, 50, 100},
			},
		},
	}

	// 创建生成器
	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// 生成代码
	opts := Options{OutputDir: outputDir}
	if err := gen.Generate(cfg, opts); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// 验证文件是否生成
	outputFile := filepath.Join(outputDir, "metrics_gen.go")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Generate() did not create file %s", outputFile)
	}

	// 读取生成的文件内容
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}
	contentStr := string(content)

	// 验证包名正确（应该是 testmetrics）
	if !strings.Contains(contentStr, "package testmetrics") {
		t.Error("Generated file does not contain correct package name")
	}

	// 验证包含关键结构
	expectedContents := []string{
		"type Metrics struct",
		"RequestsTotal",
		"*RequestsTotalMetric",
		"ActiveConnections",
		"*ActiveConnectionsMetric",
		"RequestDurationMs",
		"*RequestDurationMsMetric",
		"func NewMetrics()",
		"test-service",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Generated file missing expected content: %s\nActual content:\n%s", expected, contentStr)
		}
	}
}

func TestGetPackageName(t *testing.T) {
	tests := []struct {
		outputDir string
		expected  string
	}{
		{".", "metrics"},
		{"/", "metrics"},
		{"", "metrics"},
		{"metrics", "metrics"},
		{"aaa/bbb/ccc", "ccc"},
		{"aaa\\bbb\\ccc", "ccc"},
		{"myapp/internal/metrics", "metrics"},
		{"./metrics", "metrics"},
	}

	for _, tt := range tests {
		t.Run(tt.outputDir, func(t *testing.T) {
			result := getPackageName(tt.outputDir)
			if result != tt.expected {
				t.Errorf("getPackageName(%q) = %q, want %q", tt.outputDir, result, tt.expected)
			}
		})
	}
}

func TestGenerate_InvalidOutputDir(t *testing.T) {
	// 创建一个无法写入的目录路径（在 Windows 上使用非法字符，在 Unix 上使用 null 字节）
	var invalidDir string
	if os.PathSeparator == '\\' {
		invalidDir = "CON:" // Windows 保留设备名
	} else {
		invalidDir = "/dev/null/invalid" // Unix 上无法创建子目录
	}

	cfg := &config.MetricConfig{
		Service: "test-service",
		Metrics: []config.Metric{
			{
				Name:    "requests_total",
				Help:    "Total requests",
				Type:    "counter",
				Labels:  []config.Label{},
				Methods: []string{"inc"},
			},
		},
	}

	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	opts := Options{OutputDir: invalidDir}
	err = gen.Generate(cfg, opts)
	if err == nil {
		t.Error("Generate() with invalid output dir should return error")
	}
}

func TestGenerate_ReadOnlyOutputDir(t *testing.T) {
	// Windows 上权限测试行为不同，跳过
	if os.PathSeparator == '\\' {
		t.Skip("Skipping read-only directory test on Windows")
	}

	if os.Getuid() == 0 {
		t.Skip("Skipping test when running as root")
	}

	tmpDir := t.TempDir()
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.Mkdir(readOnlyDir, 0555); err != nil {
		t.Skipf("Cannot create read-only directory: %v", err)
	}
	defer os.Chmod(readOnlyDir, 0755) // 清理时恢复权限

	cfg := &config.MetricConfig{
		Service: "test-service",
		Metrics: []config.Metric{
			{
				Name:    "requests_total",
				Help:    "Total requests",
				Type:    "counter",
				Labels:  []config.Label{},
				Methods: []string{"inc"},
			},
		},
	}

	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	opts := Options{OutputDir: filepath.Join(readOnlyDir, "subdir")}
	err = gen.Generate(cfg, opts)
	if err == nil {
		t.Error("Generate() to read-only directory should return error")
	}
}

func TestGenerate_EmptyMetrics(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "empty")

	cfg := &config.MetricConfig{
		Service: "test-service",
		Metrics: []config.Metric{},
	}

	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	opts := Options{OutputDir: outputDir}
	if err := gen.Generate(cfg, opts); err != nil {
		t.Fatalf("Generate() with empty metrics error = %v", err)
	}

	outputFile := filepath.Join(outputDir, "metrics_gen.go")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Generate() did not create file %s", outputFile)
	}
}

func TestGenerate_MultipleMetrics(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "multi")

	cfg := &config.MetricConfig{
		Service: "test-service",
		Metrics: []config.Metric{
			{
				Name:    "requests_total",
				Help:    "Total requests",
				Type:    "counter",
				Labels:  []config.Label{{Name: "method", Vals: []string{"GET", "POST"}}},
				Methods: []string{"inc"},
			},
			{
				Name:    "connections",
				Help:    "Active connections",
				Type:    "gauge",
				Labels:  []config.Label{},
				Methods: []string{"set"},
			},
			{
				Name:    "latency_ms",
				Help:    "Latency",
				Type:    "histogram",
				Labels:  []config.Label{},
				Methods: []string{"observe"},
				Buckets: []float64{1, 10, 100},
			},
		},
	}

	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	opts := Options{OutputDir: outputDir}
	if err := gen.Generate(cfg, opts); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	outputFile := filepath.Join(outputDir, "metrics_gen.go")
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}
	contentStr := string(content)

	// 验证包含所有指标的结构
	expectedStructs := []string{
		"type Metrics struct",
		"RequestsTotal",
		"Connections",
		"LatencyMs",
	}
	for _, expected := range expectedStructs {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Generated file missing expected struct: %s", expected)
		}
	}
}

func TestGenerate_MetricWithMultipleLabels(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "labels")

	cfg := &config.MetricConfig{
		Service: "test-service",
		Metrics: []config.Metric{
			{
				Name: "requests_total",
				Help: "Total requests",
				Type: "counter",
				Labels: []config.Label{
					{Name: "method", Vals: []string{"GET", "POST"}},
					{Name: "path", Vals: []string{"*"}},
					{Name: "status", Vals: []string{"200", "404", "500"}},
				},
				Methods: []string{"inc"},
			},
		},
	}

	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	opts := Options{OutputDir: outputDir}
	if err := gen.Generate(cfg, opts); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	outputFile := filepath.Join(outputDir, "metrics_gen.go")
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}
	contentStr := string(content)

	// 验证标签被正确包含
	if !strings.Contains(contentStr, `"method"`) || !strings.Contains(contentStr, `"path"`) {
		t.Error("Generated file should contain multiple labels")
	}
}

func TestGenerate_MultipleMethods(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "methods")

	cfg := &config.MetricConfig{
		Service: "test-service",
		Metrics: []config.Metric{
			{
				Name:    "requests_total",
				Help:    "Total requests",
				Type:    "counter",
				Labels:  []config.Label{},
				Methods: []string{"inc", "add"},
			},
		},
	}

	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	opts := Options{OutputDir: outputDir}
	if err := gen.Generate(cfg, opts); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	outputFile := filepath.Join(outputDir, "metrics_gen.go")
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}
	contentStr := string(content)

	// 验证多个方法被生成
	if !strings.Contains(contentStr, "func (m *RequestsTotalMetric)") {
		t.Error("Generated file should contain counter methods")
	}
}

func TestGenerate_EnumLabels(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "enums")

	cfg := &config.MetricConfig{
		Service: "test-service",
		Metrics: []config.Metric{
			{
				Name: "requests_total",
				Help: "Total requests",
				Type: "counter",
				Labels: []config.Label{
					{Name: "source", Vals: []string{"app", "web"}},
				},
				Methods: []string{"inc"},
			},
		},
	}

	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	opts := Options{OutputDir: outputDir}
	if err := gen.Generate(cfg, opts); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	outputFile := filepath.Join(outputDir, "metrics_gen.go")
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}
	contentStr := string(content)

	// 验证枚举类型被生成
	expectedContents := []string{
		"type RequestsTotalSource interface",
		"RequestsTotalSourceApp",
		"RequestsTotalSourceWeb",
	}
	for _, expected := range expectedContents {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Generated file missing expected content: %s", expected)
		}
	}
}
