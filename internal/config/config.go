// Package config 处理 YAML 配置文件的解析
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// MetricConfig 表示 YAML 配置的根结构
type MetricConfig struct {
	ServiceName string      `yaml:"service_name"`
	Subsystems  []Subsystem `yaml:"subsystems"`
}

// Subsystem 表示一个子系统
type Subsystem struct {
	Name       string      `yaml:"name"`
	Counters   []Metric    `yaml:"counters"`
	Gauges     []Metric    `yaml:"gauges"`
	Histograms []Histogram `yaml:"histograms"`
}

// Metric 表示计数器或仪表盘指标
type Metric struct {
	Name    string   `yaml:"name"`
	Help    string   `yaml:"help"`
	Labels  []string `yaml:"labels"`
	Methods []string `yaml:"methods"`
}

// Histogram 表示直方图指标
type Histogram struct {
	Name    string    `yaml:"name"`
	Help    string    `yaml:"help"`
	Labels  []string  `yaml:"labels"`
	Methods []string  `yaml:"methods"`
	Buckets []float64 `yaml:"buckets"`
}

// Load 从指定路径加载 YAML 配置文件
func Load(path string) (*MetricConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading YAML file %s: %w", path, err)
	}

	var cfg MetricConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}

	return &cfg, nil
}
