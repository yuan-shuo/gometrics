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

	// 验证基本结构
	assertService(t, cfg, "user")
	assertMetricsCount(t, cfg, 3)

	// 验证各个指标
	assertCounterMetric(t, cfg.Metrics[0])
	assertGaugeMetric(t, cfg.Metrics[1])
	assertHistogramMetric(t, cfg.Metrics[2])
}

func assertService(t *testing.T, cfg *MetricConfig, want string) {
	t.Helper()
	if cfg.Service != want {
		t.Errorf("Service = %v, want %v", cfg.Service, want)
	}
}

func assertMetricsCount(t *testing.T, cfg *MetricConfig, want int) {
	t.Helper()
	if len(cfg.Metrics) != want {
		t.Fatalf("len(Metrics) = %v, want %v", len(cfg.Metrics), want)
	}
}

func assertCounterMetric(t *testing.T, m Metric) {
	t.Helper()
	if m.Name != "requests_total" {
		t.Errorf("Name = %v, want %v", m.Name, "requests_total")
	}
	if m.Help != "Total requests" {
		t.Errorf("Help = %v, want %v", m.Help, "Total requests")
	}
	if m.Type != "counter" {
		t.Errorf("Type = %v, want %v", m.Type, "counter")
	}
	if len(m.Labels) != 2 {
		t.Errorf("Labels len = %v, want %v", len(m.Labels), 2)
	}
	if m.Labels[0].Name != "method" {
		t.Errorf("Labels[0].Name = %v, want %v", m.Labels[0].Name, "method")
	}
	if !sliceEqual(m.Labels[0].Vals, []string{"GET", "POST"}) {
		t.Errorf("Labels[0].Vals = %v, want %v", m.Labels[0].Vals, []string{"GET", "POST"})
	}
	if !sliceEqual(m.Methods, []string{"inc", "add"}) {
		t.Errorf("Methods = %v, want %v", m.Methods, []string{"inc", "add"})
	}
}

func assertGaugeMetric(t *testing.T, m Metric) {
	t.Helper()
	if m.Name != "active_connections" {
		t.Errorf("Name = %v, want %v", m.Name, "active_connections")
	}
	if m.Type != "gauge" {
		t.Errorf("Type = %v, want %v", m.Type, "gauge")
	}
}

func assertHistogramMetric(t *testing.T, m Metric) {
	t.Helper()
	if m.Name != "request_duration_ms" {
		t.Errorf("Name = %v, want %v", m.Name, "request_duration_ms")
	}
	if m.Type != "histogram" {
		t.Errorf("Type = %v, want %v", m.Type, "histogram")
	}
	if !floatSliceEqual(m.Buckets, []float64{10, 50, 100}) {
		t.Errorf("Buckets = %v, want %v", m.Buckets, []float64{10, 50, 100})
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
