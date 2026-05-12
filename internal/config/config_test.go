package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	// 创建临时测试文件
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.yaml")

	yamlContent := `
service: user
metrics:
  - name: requests_total
    help: Total requests
    type: counter
    labels:
      - name: method
        vals: ["GET", "POST"]
      - name: path
        vals: ["*"]
    methods: ["inc", "add"]
  - name: active_connections
    help: Active connections
    type: gauge
    labels:
      - name: pool
        vals: ["default", "cache"]
    methods: ["set", "inc", "dec"]
  - name: request_duration_ms
    help: Request duration
    type: histogram
    labels:
      - name: method
        vals: ["GET", "POST"]
    buckets: [10, 50, 100]
    methods: ["observe"]
`
	if err := os.WriteFile(testFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// 测试加载配置
	cfg, err := Load(testFile)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// 验证 service
	if cfg.Service != "user" {
		t.Errorf("Service = %v, want %v", cfg.Service, "user")
	}

	// 验证 metrics 数量
	if len(cfg.Metrics) != 3 {
		t.Fatalf("len(Metrics) = %v, want %v", len(cfg.Metrics), 3)
	}

	// 验证第一个 counter 指标
	counter := cfg.Metrics[0]
	if counter.Name != "requests_total" {
		t.Errorf("Metric[0].Name = %v, want %v", counter.Name, "requests_total")
	}
	if counter.Help != "Total requests" {
		t.Errorf("Metric[0].Help = %v, want %v", counter.Help, "Total requests")
	}
	if counter.Type != "counter" {
		t.Errorf("Metric[0].Type = %v, want %v", counter.Type, "counter")
	}
	if len(counter.Labels) != 2 {
		t.Errorf("Metric[0].Labels len = %v, want %v", len(counter.Labels), 2)
	}
	if counter.Labels[0].Name != "method" {
		t.Errorf("Metric[0].Labels[0].Name = %v, want %v", counter.Labels[0].Name, "method")
	}
	if !sliceEqual(counter.Labels[0].Vals, []string{"GET", "POST"}) {
		t.Errorf("Metric[0].Labels[0].Vals = %v, want %v", counter.Labels[0].Vals, []string{"GET", "POST"})
	}
	if !sliceEqual(counter.Methods, []string{"inc", "add"}) {
		t.Errorf("Metric[0].Methods = %v, want %v", counter.Methods, []string{"inc", "add"})
	}

	// 验证 gauge 指标
	gauge := cfg.Metrics[1]
	if gauge.Name != "active_connections" {
		t.Errorf("Metric[1].Name = %v, want %v", gauge.Name, "active_connections")
	}
	if gauge.Type != "gauge" {
		t.Errorf("Metric[1].Type = %v, want %v", gauge.Type, "gauge")
	}

	// 验证 histogram 指标
	histogram := cfg.Metrics[2]
	if histogram.Name != "request_duration_ms" {
		t.Errorf("Metric[2].Name = %v, want %v", histogram.Name, "request_duration_ms")
	}
	if histogram.Type != "histogram" {
		t.Errorf("Metric[2].Type = %v, want %v", histogram.Type, "histogram")
	}
	if !floatSliceEqual(histogram.Buckets, []float64{10, 50, 100}) {
		t.Errorf("Metric[2].Buckets = %v, want %v", histogram.Buckets, []float64{10, 50, 100})
	}
}

func TestLoad_FileNotExist(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("Load() expected error for non-existent file, got nil")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "invalid.yaml")

	// 写入无效的 YAML
	if err := os.WriteFile(testFile, []byte("invalid: yaml: content: ["), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	_, err := Load(testFile)
	if err == nil {
		t.Error("Load() expected error for invalid YAML, got nil")
	}
}

func TestMetric_GetLabelNames(t *testing.T) {
	metric := Metric{
		Name: "test_metric",
		Labels: []Label{
			{Name: "label1", Vals: []string{"a", "b"}},
			{Name: "label2", Vals: []string{"*"}},
			{Name: "label3", Vals: []string{"x", "y", "z"}},
		},
	}

	names := metric.GetLabelNames()
	expected := []string{"label1", "label2", "label3"}
	if !sliceEqual(names, expected) {
		t.Errorf("GetLabelNames() = %v, want %v", names, expected)
	}
}

func TestMetric_GetLabelNames_Empty(t *testing.T) {
	metric := Metric{
		Name:   "test_metric",
		Labels: []Label{},
	}

	names := metric.GetLabelNames()
	if len(names) != 0 {
		t.Errorf("GetLabelNames() = %v, want empty slice", names)
	}
}

func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func floatSliceEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
