# gometrics

一个用于生成 Prometheus 指标代码的 Go 代码生成工具。通过 YAML 配置文件定义指标，自动生成类型安全的指标管理代码。

## 特性

- 支持 Counter、Gauge、Histogram 三种指标类型
- 类型安全的标签参数
- 自动根据输出目录确定包名
- 基于 Go 模板生成代码
- 支持自定义方法（Inc、Add、Set、Dec、Observe）

## 安装

```bash
go install github.com/yuan-shuo/gometrics@latest
```

或者从源码构建：

```bash
git clone https://github.com/yuan-shuo/gometrics.git
cd gometrics
go build -o gometrics .
```

## 使用方法

### 1. 创建 YAML 配置文件

参考 `metrics.example.yaml` 创建你的配置文件：

```yaml
service_name: "my-service"
subsystems:
  - name: "order"
    counters:
      - name: "orders_created_total"
        help: "Total number of orders created"
        labels: ["status", "channel"]
        methods: ["inc"]
    gauges:
      - name: "active_connections"
        help: "Current active connections"
        labels: ["pool"]
        methods: ["set", "inc", "dec"]
    histograms:
      - name: "request_latency_ms"
        help: "Request latency in milliseconds"
        labels: ["method", "path"]
        buckets: [10, 50, 100, 200, 500, 1000, 2000]
        methods: ["observe"]
```

### 2. 运行代码生成工具

```bash
# 生成代码到指定目录
gometrics -f metrics.yaml -d ./internal/metrics

# 示例：生成到 metrics 目录
gometrics -f metrics.yaml -d metrics
```

### 3. 在项目中使用生成的代码

```go
package main

import (
    "your-project/metrics"
)

func main() {
    // 创建指标管理器
    m := metrics.NewMetricsManager()
    
    // 使用 Counter
    m.Order.OrdersCreatedTotal.Inc("success", "web")
    
    // 使用 Gauge
    m.Order.ActiveConnections.Set(100, "pool1")
    m.Order.ActiveConnections.Inc("pool1")
    m.Order.ActiveConnections.Dec("pool1")
    
    // 使用 Histogram
    m.Order.RequestLatencyMs.Observe(150, "GET", "/api/orders")
}
```

## 配置文件说明

### 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| `service_name` | string | 服务名称，作为指标的 namespace |
| `subsystems` | array | 子系统列表 |
| `subsystems[].name` | string | 子系统名称，作为指标的 subsystem |
| `subsystems[].counters` | array | 计数器指标列表 |
| `subsystems[].gauges` | array | 仪表盘指标列表 |
| `subsystems[].histograms` | array | 直方图指标列表 |

### 指标字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `name` | string | 指标名称（snake_case） |
| `help` | string | 指标帮助信息 |
| `labels` | array | 标签列表 |
| `methods` | array | 支持的方法 |
| `buckets` | array | 直方图桶边界（仅 histogram） |

### 支持的方法

| 指标类型 | 方法 | 说明 |
|----------|------|------|
| Counter | `inc` | 递增计数器 |
| Counter | `add` | 增加指定值 |
| Gauge | `set` | 设置值 |
| Gauge | `inc` | 递增 |
| Gauge | `dec` | 递减 |
| Histogram | `observe` | 记录观察值 |

## 命令行参数

```
Usage of gometrics:
  -d string
        Output directory for the generated Go file (required)
  -f string
        Path to the YAML configuration file (required)
```

## 包名规则

生成的代码包名根据输出目录自动确定：

| 输出目录 | 包名 |
|----------|------|
| `.` | `metrics` |
| `metrics` | `metrics` |
| `aaa/bbb/ccc` | `ccc` |

## 项目结构

```
.
├── main.go                      # 入口文件
├── internal/
│   ├── config/
│   │   └── config.go           # YAML 配置解析
│   ├── template/
│   │   └── template.go         # 代码模板
│   └── generator/
│       └── generator.go        # 代码生成逻辑
├── metrics.example.yaml         # 示例配置文件
└── README.md                    # 本文档
```

## 示例

查看 `metrics.example.yaml` 获取完整的配置示例。

生成示例代码：

```bash
gometrics -f metrics.example.yaml -d ./example
```