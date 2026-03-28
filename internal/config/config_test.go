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
service_name: "test-service"
subsystems:
  - name: "api"
    counters:
      - name: "requests_total"
        help: "Total requests"
        labels: ["method", "path"]
        methods: ["inc"]
    gauges:
      - name: "active_connections"
        help: "Active connections"
        labels: ["pool"]
        methods: ["set", "inc", "dec"]
    histograms:
      - name: "request_duration_ms"
        help: "Request duration"
        labels: ["method"]
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
	assertServiceName(t, cfg, "test-service")
	assertSubsystemCount(t, cfg, 1)

	// 验证子系统
	subsystem := cfg.Subsystems[0]
	assertSubsystemName(t, subsystem, "api")

	// 验证 counters
	assertCounterCount(t, subsystem, 1)
	assertCounter(t, subsystem.Counters[0], "requests_total", "Total requests", []string{"method", "path"}, []string{"inc"})

	// 验证 gauges
	assertGaugeCount(t, subsystem, 1)
	assertGaugeName(t, subsystem.Gauges[0], "active_connections")

	// 验证 histograms
	assertHistogramCount(t, subsystem, 1)
	assertHistogram(t, subsystem.Histograms[0], "request_duration_ms", []float64{10, 50, 100})
}

func assertServiceName(t *testing.T, cfg *MetricConfig, want string) {
	t.Helper()
	if cfg.ServiceName != want {
		t.Errorf("ServiceName = %v, want %v", cfg.ServiceName, want)
	}
}

func assertSubsystemCount(t *testing.T, cfg *MetricConfig, want int) {
	t.Helper()
	if len(cfg.Subsystems) != want {
		t.Fatalf("len(Subsystems) = %v, want %v", len(cfg.Subsystems), want)
	}
}

func assertSubsystemName(t *testing.T, s Subsystem, want string) {
	t.Helper()
	if s.Name != want {
		t.Errorf("Subsystem.Name = %v, want %v", s.Name, want)
	}
}

func assertCounterCount(t *testing.T, s Subsystem, want int) {
	t.Helper()
	if len(s.Counters) != want {
		t.Fatalf("len(Counters) = %v, want %v", len(s.Counters), want)
	}
}

func assertCounter(t *testing.T, c Metric, name, help string, labels, methods []string) {
	t.Helper()
	if c.Name != name {
		t.Errorf("Counter.Name = %v, want %v", c.Name, name)
	}
	if c.Help != help {
		t.Errorf("Counter.Help = %v, want %v", c.Help, help)
	}
	if !sliceEqual(c.Labels, labels) {
		t.Errorf("Counter.Labels = %v, want %v", c.Labels, labels)
	}
	if !sliceEqual(c.Methods, methods) {
		t.Errorf("Counter.Methods = %v, want %v", c.Methods, methods)
	}
}

func assertGaugeCount(t *testing.T, s Subsystem, want int) {
	t.Helper()
	if len(s.Gauges) != want {
		t.Fatalf("len(Gauges) = %v, want %v", len(s.Gauges), want)
	}
}

func assertGaugeName(t *testing.T, g Metric, want string) {
	t.Helper()
	if g.Name != want {
		t.Errorf("Gauge.Name = %v, want %v", g.Name, want)
	}
}

func assertHistogramCount(t *testing.T, s Subsystem, want int) {
	t.Helper()
	if len(s.Histograms) != want {
		t.Fatalf("len(Histograms) = %v, want %v", len(s.Histograms), want)
	}
}

func assertHistogram(t *testing.T, h Histogram, name string, buckets []float64) {
	t.Helper()
	if h.Name != name {
		t.Errorf("Histogram.Name = %v, want %v", h.Name, name)
	}
	if !floatSliceEqual(h.Buckets, buckets) {
		t.Errorf("Histogram.Buckets = %v, want %v", h.Buckets, buckets)
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
