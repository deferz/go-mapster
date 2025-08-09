# go-mapster 新架构公开 API 设计草案 (v1 Draft)

> 目标：先冻结「对外 API 面」再逐步实现。所有实现以 *Plan 驱动 + 可插拔扩展* 为核心。本文将区分 **M0 最小集** 与 **后续阶段能力**，避免一次性过度设计。

---
## 目录
1. 顶层功能分层
2. 最小可用集合 (M0) API
3. 渐进能力 (M1~M5)
4. 配置 Builder 设计
5. 运行期 Options 设计
6. Converter / Hook / Plan 生态接口
7. 错误模型
8. Codegen 预留接口
9. 版本与兼容策略
10. Open Questions

---
## 1. 顶层功能分层
```
User API Layer
  |-- Mapping Functions (Map / MapE / MapTo / MapSlice*)
  |-- Configuration DSL (Config[S,T])
  |-- Options (WithStrict / WithContext / ...)
  |-- Registry Ops (ListPlans / Invalidate / Stats)
Extension Layer
  |-- Converters (type-pair transformations)
  |-- Hooks (pre/post execution)
  |-- Tag Strategies (struct tag resolving)
  |-- Codegen (Plan -> Source)
Core Layer
  |-- Plan Builder (analyze + compile field ops)
  |-- Plan Cache (immutable, versioned)
  |-- Executor (pipeline: pre-hooks -> ops -> post-hooks)
  |-- Introspection (expose plan / field meta)
Infrastructure
  |-- Thread-safe registries
  |-- Metrics & Tracing interfaces (optional)
  |-- Error taxonomy
```

---
## 2. 最小可用集合 (M0) API
仅支持：struct→struct，字段同名/重命名映射，条件，单步 transform，忽略字段。

### 核心函数
```go
// 忽略错误（失败返回零值）
func Map[T any](src any) T

// 显式错误返回
func MapE[T any](src any, opts ...Option) (T, error)

// 就地映射
func MapTo[T any](src any, target *T, opts ...Option) error
```

### 配置 DSL (M0 子集)
```go
func Config[S, T any]() *ConfigBuilder[S,T]

// 链式：
(builder).Map("TargetField").From("SourceField").Transform(fn).When(cond).Done()
(builder).Ignore("Field")
(builder).Register() error
```

### Transform 支持签名
```go
func(any) any
func(any) (any, error)
func(SrcFieldType) OutType
func(SrcFieldType) (OutType, error)
```

### Condition 签名
```go
func(S) bool  // S 为整个源对象
```

### Option (M0)
```go
WithStrict() // 严格模式：任何字段级错误直接返回
```

---
## 3. 渐进能力 (Roadmap)
| 阶段 | 能力 | 说明 |
|------|------|------|
| M1 | Converter Registry / FromPath 预编译 / 错误分类 | 提供类型对转换 & 路径缓存 & FieldError 分类 |
| M2 | Slice / Array / Map 容器映射 + Context | MapSlice, MapSliceE, 上下文透传 Hooks |
| M3 | Auto Flatten 集成 Plan / 多 Transform 链 / 循环检测 Hook | Flatten 提前编译；`TransformChain` 支持 |
| M4 | Metrics & Tracing / TagStrategy(JSON等) / Codegen IR | 生成器 CLI & 热点统计 |
| M5 | 并发流水线优化 / Unsafe 快速路径(可选 build tag) | 性能冲刺 |

---
## 4. 配置 Builder 设计
```go
type ConfigBuilder[S,T any] struct { /* internal */ }

func (b *ConfigBuilder[S,T]) Map(target string) *FieldBuilder[S,T]
func (b *ConfigBuilder[S,T]) Ignore(target string) *ConfigBuilder[S,T]
func (b *ConfigBuilder[S,T]) Register() error

// Field Builder
 type FieldBuilder[S,T any] struct { /* internal */ }
 func (f *FieldBuilder[S,T]) From(source string) *FieldBuilder[S,T]
 func (f *FieldBuilder[S,T]) FromPath(path string) *FieldBuilder[S,T]            // M1
 func (f *FieldBuilder[S,T]) When(cond func(S) bool) *FieldBuilder[S,T]
 func (f *FieldBuilder[S,T]) Transform(fn any) *FieldBuilder[S,T]                // M0 单步
 func (f *FieldBuilder[S,T]) TransformChain(fns ...any) *FieldBuilder[S,T]       // M3
 func (f *FieldBuilder[S,T]) ConvertWith(converterID string) *FieldBuilder[S,T]  // M1 绑定注册转换
 func (f *FieldBuilder[S,T]) Flatten(depth int) *FieldBuilder[S,T]               // M3 针对嵌套聚合
 func (f *FieldBuilder[S,T]) Done() *ConfigBuilder[S,T]
```

### 注册后行为
- 生成 `MappingDefinition` 放入 registry.defs
- Plan 首次命中构建：解析 FieldRules -> 编译成 FieldOps
- FromPath 在构建时解析成 accessor 链（结构数组）

---
## 5. 运行期 Options 设计
```go
type Option interface { apply(*mapOptions) }

type mapOptions struct {
    strict        bool
    ctx           context.Context       // M2
    maxDepth      int                   // M3 (循环/flatten 限制)
    enableFlatten bool                  // M3
    converters    []ConverterOverride   // M1 指定优先 converter
    tracer        Tracer                // M4
    metrics       MetricsCollector      // M4
    poolPolicy    PoolPolicy            // M5
    partial       bool                  // M2 允许部分成功 + 错误收集
    errorSink     ErrorCollector        // M2/M3
}

// 已计划 Option 构造函数
WithStrict()
WithContext(ctx context.Context)            // M2
WithMaxDepth(n int)                         // M3
WithFlatten(enabled bool)                   // M3
WithConverters(ids ...string)               // M1
WithTracer(t Tracer)                        // M4
WithMetrics(m MetricsCollector)             // M4
WithObjectPool(policy PoolPolicy)           // M5
WithPartial()                               // M2
WithErrorCollector(c ErrorCollector)        // M2
```

---
## 6. Converter / Hook / Plan 生态接口
### Converter
```go
type Converter interface {
    ID() string
    Source() reflect.Type      // or generic registration helper
    Target() reflect.Type
    Convert(input any) (any, error)
    Cost() int                 // 选择优先级，数值越低优先
    Lossy() bool               // 是否可能信息丢失（规划）
}

// 注册
func RegisterConverter[S,T any](id string, fn func(S) (T,error), opts ...ConverterOption) error // M1

// 查找策略
// 1. 精确 (S->T)
// 2. 可逆或链式（未来可选）
```

### Hook (Pipeline)
```go
type Hook interface {
    Before(plan *Plan, ctx HookContext) error
    After(plan *Plan, ctx HookContext) error
}

type HookContext struct {
    Context context.Context
    Src     any
    Dst     any // pointer during After
    Options mapOptions
    Errors  []error            // 部分模式
    Meta    map[string]any
}

func RegisterHook(h Hook) // 注册全局 Hook (内置顺序：metrics -> tracing -> user -> debug)
```

### Plan Introspection
```go
type PlanInfo struct {
    Src, Dst    reflect.Type
    FieldOps    []FieldOpInfo
    Flags       PlanFlags
    Version     uint32
}

type FieldOpInfo struct {
    TargetField string
    SourcePath  string
    ConverterID string
    HasCondition bool
    HasTransform bool
}

func DescribePlan[S,T any]() (*PlanInfo, bool)
func ListPlans() []PlanInfo
func Invalidate[S,T any]() bool
func InvalidateAll()
```

---
## 7. 错误模型
```go
type MappingError interface { error; Unwrap() error; Is(target error) bool }

var (
    ErrPlanBuild        = errors.New("plan build failed")
    ErrInvalidSignature = errors.New("invalid function signature")
    ErrTypeMismatch     = errors.New("type mismatch")
    ErrConverterMissing = errors.New("converter not found")
    ErrMaxDepthExceeded = errors.New("max depth exceeded")
)

type FieldError struct {
    Field      string
    Reason     error
    SourcePath string
}

type AggregateError struct { Fields []FieldError }
```
- *宽松模式*：字段错误 -> 收集 (若提供 ErrorCollector)，否则静默。
- *严格模式*：首个错误直接返回。
- *部分模式* (WithPartial)：返回目标对象 + AggregateError。

---
## 8. Codegen 预留接口 (M4)
```go
type CodegenRequest struct {
    Plans []PlanInfo
    OutputDir string
    Package   string
    Strategy  CodegenStrategy // 单文件 / 每 Plan 分文件
}

type CodegenResult struct {
    Files []string
    Warnings []string
}

func GenerateCode(req CodegenRequest) (CodegenResult, error)
```
生成文件包含：
- `func map_<hash>(src <Src>) <Dst>`
- `func init(){ mapster.RegisterGeneratedMapper(map_<hash>) }`

---
## 9. 版本与兼容策略
- v0.x：快速迭代，不保证二进制兼容（主分支）。
- v1.0：稳定核心 API：Map/MapE/MapTo/Config DSL/Options基础/Converter/Hook。
- v2 (可选模块路径 `/v2`)：引入重大行为变化（如上下文强制、默认严格等）。

---
## 10. Open Questions (需要与使用者进一步确认)
| 主题 | 问题 | 需要输入 |
|------|------|----------|
| Tag 支持 | 解析 json / db / mapster tag 的优先级? | 期望字段来源规则 |
| 并发策略 | Plan 构建是否允许阻塞还是采用乐观重试? | 高并发场景指标 |
| Codegen | 是否需要根据热点调用自动输出候选? | 使用频率感知阈值 |
| Unsafe | 是否允许通过 build tag 启用？ | 运行环境限制 |
| Converter 链 | 是否允许 A->B->C 自动推导? | 复杂性 vs 价值 |
| 循环检测 | 默认启用还是按需? | 使用场景频率 |

---
## 结论
以上 API 设计确保：
1. 最小集足以替换基础手写映射样板。
2. 渐进路线避免一次实现所有复杂特性。
3. 足够的扩展点（Converter / Hook / Codegen / TagStrategy）。

> 请审阅：若你认可 M0 API 与整体扩展骨架，我将：
> 1. 清空旧实现文件（若还未执行）
> 2. 提交 M0 骨架代码 (registry + plan + executor + config builder + tests)
> 3. 开始 M1 转换器与路径预编译实现。

反馈点欢迎直接在本文件追加 TODO，我们会版本号标注。
