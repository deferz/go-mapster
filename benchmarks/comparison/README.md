# Go 对象映射库性能对比

这是一个独立的性能测试项目，用于对比 Go 生态中常见的对象映射库。

## 测试的库

1. **[mapster](https://github.com/deferz/go-mapster)** - 本项目，高性能对象映射库
2. **[copier](https://github.com/jinzhu/copier)** - 流行的结构体复制库
3. **[mergo](https://github.com/imdario/mergo)** - 合并和转换结构体的库
4. **[mapper](https://github.com/devfeel/mapper)** - 简单快速的对象映射器
5. **[mapstructure](https://github.com/mitchellh/mapstructure)** - 用于解码通用 map 到结构体的库

## 运行测试

```bash
# 下载依赖
go mod download

# 运行所有基准测试
go test -bench=. -benchmem

# 运行特定测试
go test -bench=BenchmarkSimpleMapping -benchmem

# 生成性能分析图表
go test -bench=. -benchmem -cpuprofile=cpu.prof
go tool pprof -http=:8080 cpu.prof
```

## 测试场景

1. **简单对象映射** - 基本字段复制
2. **复杂对象映射** - 包含嵌套结构和自定义逻辑
3. **批量映射** - 切片和数组的映射
4. **大对象映射** - 多字段的复杂结构体

## 注意事项

- 这是一个独立的模块，不会影响主项目的依赖
- 测试结果可能因硬件和 Go 版本而异
- 各库的功能特性不完全相同，性能只是选择的一个方面 