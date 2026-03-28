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

	// 验证 service_name
	if cfg.ServiceName != "test-service" {
		t.Errorf("ServiceName = %v, want %v", cfg.ServiceName, "test-service")
	}

	// 验证子系统
	if len(cfg.Subsystems) != 1 {
		t.Fatalf("len(Subsystems) = %v, want %v", len(cfg.Subsystems), 1)
	}

	subsystem := cfg.Subsystems[0]
	if subsystem.Name != "api" {
		t.Errorf("Subsystem.Name = %v, want %v", subsystem.Name, "api")
	}

	// 验证 counters
	if len(subsystem.Counters) != 1 {
		t.Fatalf("len(Counters) = %v, want %v", len(subsystem.Counters), 1)
	}
	counter := subsystem.Counters[0]
	if counter.Name != "requests_total" {
		t.Errorf("Counter.Name = %v, want %v", counter.Name, "requests_total")
	}
	if counter.Help != "Total requests" {
		t.Errorf("Counter.Help = %v, want %v", counter.Help, "Total requests")
	}
	if len(counter.Labels) != 2 || counter.Labels[0] != "method" || counter.Labels[1] != "path" {
		t.Errorf("Counter.Labels = %v, want [method path]", counter.Labels)
	}
	if len(counter.Methods) != 1 || counter.Methods[0] != "inc" {
		t.Errorf("Counter.Methods = %v, want [inc]", counter.Methods)
	}

	// 验证 gauges
	if len(subsystem.Gauges) != 1 {
		t.Fatalf("len(Gauges) = %v, want %v", len(subsystem.Gauges), 1)
	}
	gauge := subsystem.Gauges[0]
	if gauge.Name != "active_connections" {
		t.Errorf("Gauge.Name = %v, want %v", gauge.Name, "active_connections")
	}

	// 验证 histograms
	if len(subsystem.Histograms) != 1 {
		t.Fatalf("len(Histograms) = %v, want %v", len(subsystem.Histograms), 1)
	}
	histogram := subsystem.Histograms[0]
	if histogram.Name != "request_duration_ms" {
		t.Errorf("Histogram.Name = %v, want %v", histogram.Name, "request_duration_ms")
	}
	if len(histogram.Buckets) != 3 || histogram.Buckets[0] != 10 {
		t.Errorf("Histogram.Buckets = %v, want [10 50 100]", histogram.Buckets)
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
