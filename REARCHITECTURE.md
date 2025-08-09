# Go Mapster 重构设计（Draft v0.1）

> 目标：在保持现有使用体验的同时，解决当前架构在可维护性、扩展性、并发安全、性能可预见性方面的问题，为未来自动代码生成与插件化能力打基础。

## 1. 现存主要问题（基于代码审阅 + 合理假设）

| 类别 | 问题 | 影响 |
|------|------|------|
| 并发安全 | 全局 `generatedMappers` / `globalConfigs` 未加锁 | 运行期并发注册/映射会产生数据竞争 |
| 配置生命周期 | 注册后不可撤销/隔离，无法按上下文（租户/请求）定制 | 不利于多租户 / 测试隔离 |
| 扩展性 | 时间转换、扁平化、路径解析等为散落逻辑，缺统一扩展点 | 新增特性需修改核心代码，风险高 |
| 性能稳定性 | 反射路径每次遍历字段，无计划缓存“映射计划” | 高频调用抖动、GC 压力增加 |
| 错误处理 | 静默失败（类型不兼容/函数签名错误直接忽略） | 难以排查错误、生产风险高 |
| 映射链路 | 逻辑分散：`mapReflect` / `mapWithConfig` / 条件 / transform | 增加修改成本、难做阶段性优化 |
| 自动扁平化 | 功能独立但尚未与映射管线融合（冲突策略/缓存缺失） | 重复解析、潜在性能损耗 |
| 转换机制 | 仅支持字段级定制函数，不支持“类型对”级转换（Converter） | 复用差，样板代码多 |
| Slice 映射 | 简单 for-loop，缺乏并行 / 预分配优化策略 | 大批量映射吞吐不足 |
| 代码生成 | 仅手动注册，缺 go:generate / 插件描述 | 零反射能力不可扩展、门槛高 |
| 可观测性 | 缺基准点埋点/统计（命中率/路径） | 难以量化优化效果 |

## 2. 重构总体目标

1. 提供 **Pipeline 化映射执行模型**（解析阶段 -> 绑定阶段 -> 运行阶段）。
2. 引入 **类型对(TypePair) -> 映射计划(Plan)** 的编译缓存，避免重复反射。
3. 抽象 **Converter Registry**（类型对转换）与 **Field Resolver / Path Resolver**。
4. 统一 **Options**（每次调用的可变行为）与 **Config**（静态注册）。
5. 提供 **并发安全的 Registry**，支持懒加载和按需失效。
6. 设计 **可插拔阶段（Hook / Middleware）**：前置校验、循环检测、value transform、post finalize。
7. 规划 **代码生成接口**：从 Plan 导出生成源码 + 运行期回填。
8. 增强 **错误与调试模式**：严格模式（panic / error）与宽松模式（跳过）。
9. 为 v1 保持兼容，v2 以 `module /v2` 形式引入改进 API。

## 3. 新架构概览

```text
Caller
  |-- Map[T]() / MapWithOptions[T]()
          |-- Lookup Plan (TypePair)
                |-- (cache miss) -> Plan Builder
                        |-- Collect Fields
                        |-- Apply Config (field rules, ignore, path, transform)
                        |-- Resolve Flatten (if enabled)
                        |-- Attach Converters
                        |-- Freeze Plan (immutable)
          |-- Execute Pipeline
                |-- PreHooks
                |-- Field Loop (fast path via pre-bound accessors)
                |-- Apply Converters / Transforms
                |-- PostHooks
                |-- Return result / error
```

### 关键组件
- `registry`：并发安全存储（generated mappers / plans / converters / configs）。
- `plan`：结构化表示一个类型对的映射策略（字段数组 + 操作序列 + flags）。
- `converter`：注册类型对转换函数 `func(S) (T, error)`；支持优先级与可逆性标记。
- `pipeline`：执行器；内部可插入 middleware。
- `options`：调用层行为开关，如：StrictMode、EnableFlatten、MaxDepth override、PoolPolicy。
- `codegen`：根据 Plan 输出模板（未来扩展）。

## 4. 数据结构草案

```go
// TypePair 唯一标识
type TypePair struct { Src, Dst reflect.Type }

// FieldOp 描述单字段操作
type FieldOp struct {
    TargetIndex int              // 目标字段索引
    SourceAccess Accessor        // 预编译的取值器
    ConverterID  uint32          // 若需类型转换
    Transform    TransformFunc   // 字段局部 transform（链式）
    Condition    ConditionFunc   // 条件函数（可选）
    Flags        FieldFlags
}

// Plan 映射计划（只读，支持代码生成）
type Plan struct {
    Pair        TypePair
    FieldOps    []FieldOp
    Meta        PlanMeta // 包含 flatten/pipeline 标志
    Version     uint32
}
```

## 5. Pipeline Hook 设计

```go
type Hook interface {
    Before(plan *Plan, ctx *Context) error
    After(plan *Plan, ctx *Context) error
}
```
- 内置 Hook：循环检测、统计埋点、调试日志。
- 用户可注册自定义 Hook。

## 6. 并发与缓存策略
- 使用 `sync.RWMutex` + `sync.Map` 混合：计划缓存 `sync.Map`，注册结构使用 RWMutex 保证一致性。
- 计划构建采用双重检查（DCL）避免重复构建。
- 支持 `Invalidate(TypePair)` 与 `InvalidateAll()`。

## 7. 错误策略
- 默认宽松（非致命失败忽略并统计）。
- `StrictMode`：类型不兼容、函数签名错误直接返回 error。
- 计划阶段进行大量预验证，运行期尽量无反射错误分支。

## 8. 代码生成规划
- `Plan -> IR`（中间结构） -> 模板生成 Go 文件 -> 用户 `go generate` 触发。
- 运行后自动调用 `RegisterGeneratedMapper` 回填。
- 未来可支持：检测热点（统计调用次数 > 阈值） → 输出待生成建议。

## 9. 迁移策略（非破坏）
1. Phase 1：引入新内部结构（当前提交）。
2. Phase 2：用 Plan + Pipeline 重写现有反射路径，老 API 仍调用新实现。
3. Phase 3：实现 Converter Registry & Options。
4. Phase 4：接入自动扁平化进入 Plan 构建阶段。
5. Phase 5：引入 Hook / 统计。
6. Phase 6：生成器 CLI / go:generate。
7. Phase 7：发布 v2（语义化错误、上下文参数）。

## 10. 初始任务拆分 (Backlog Draft)
- [ ] registry: 并发安全映射计划/转换器存储
- [ ] plan: 构建器（扫描结构 + 应用配置）
- [ ] pipeline: 执行入口（暂时包装现有 reflection 逻辑）
- [ ] converter: API & 注册验证
- [ ] options: Option 模式 + 调用包装
- [ ] metrics: 预留接口
- [ ] codegen: stub + 接口定义

## 11. API 拓展（示例）
```go
res, err := mapster.MapWithOptions[Dst](src, mapster.WithStrict(), mapster.WithFlattenDepth(2))
```

保持原：`Map[T](src)` → 内部调用 `MapWithOptions`（忽略 error 返回零值）。

## 12. 后续讨论点
- 是否需要 context.Context 透传（取消/trace）。
- 是否支持对象池（减少 slice 分配）。
- 是否引入 unsafe 快速路径（可选编译标签）。
- 是否支持 JSON tag 映射（可配置解析）。

---
此文档为初稿，欢迎补充你在生产中遇到的具体痛点，我会据此细化 Plan Builder 与错误策略。
