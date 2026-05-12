// Package config 处理 YAML 配置文件的解析
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// MetricConfig 表示 YAML 配置的根结构
type MetricConfig struct {
	Service string   `yaml:"service"`
	Metrics []Metric `yaml:"metrics"`
}

// Label 表示指标标签
type Label struct {
	Name string   `yaml:"name"`
	Vals []string `yaml:"vals"`
}

// Metric 表示单个指标配置
type Metric struct {
	Name    string  `yaml:"name"`
	Help    string  `yaml:"help"`
	Type    string  `yaml:"type"`
	Labels  []Label `yaml:"labels"`
	Methods []string `yaml:"methods"`
	Buckets []float64 `yaml:"buckets"`
}

// GetLabelNames 返回标签名称列表
func (m *Metric) GetLabelNames() []string {
	names := make([]string, len(m.Labels))
	for i, l := range m.Labels {
		names[i] = l.Name
	}
	return names
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
