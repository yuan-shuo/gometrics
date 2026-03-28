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

	// 创建测试配置
	cfg := &config.MetricConfig{
		ServiceName: "test-service",
		Subsystems: []config.Subsystem{
			{
				Name: "api",
				Counters: []config.Metric{
					{
						Name:    "requests_total",
						Help:    "Total requests",
						Labels:  []string{"method"},
						Methods: []string{"inc"},
					},
				},
				Gauges: []config.Metric{
					{
						Name:    "active_connections",
						Help:    "Active connections",
						Labels:  []string{"pool"},
						Methods: []string{"set"},
					},
				},
				Histograms: []config.Histogram{
					{
						Name:    "request_duration_ms",
						Help:    "Request duration",
						Labels:  []string{"path"},
						Methods: []string{"observe"},
						Buckets: []float64{10, 50, 100},
					},
				},
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

	// 验证包含关键结构（使用制表符而非空格，因为生成的代码使用制表符缩进）
	expectedContents := []string{
		"type ApiMetrics struct",
		"RequestsTotal",
		"*SafeCounter_RequestsTotal",
		"ActiveConnections",
		"*SafeGauge_ActiveConnections",
		"RequestDurationMs",
		"*SafeHistogram_RequestDurationMs",
		"func NewMetricsManager()",
		"namespace:",
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
