// gometrics 是一个代码生成工具，根据 YAML 配置生成 Prometheus 指标代码
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yuan-shuo/gometrics/internal/config"
	"github.com/yuan-shuo/gometrics/internal/generator"
)

func main() {
	// 解析命令行参数
	var (
		yamlFile  = flag.String("f", "", "Path to the YAML configuration file (required)")
		outputDir = flag.String("d", "", "Output directory for the generated Go file (required)")
	)
	flag.Parse()

	// 验证必选参数
	if *yamlFile == "" {
		fmt.Fprintf(os.Stderr, "Error: -f is required\n")
		flag.Usage()
		os.Exit(1)
	}
	if *outputDir == "" {
		fmt.Fprintf(os.Stderr, "Error: -d is required\n")
		flag.Usage()
		os.Exit(1)
	}

	// 加载配置
	cfg, err := config.Load(*yamlFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// 创建生成器
	gen, err := generator.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// 生成代码
	opts := generator.Options{
		OutputDir: *outputDir,
	}
	if err := gen.Generate(cfg, opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s/metrics_gen.go\n", *outputDir)
}
