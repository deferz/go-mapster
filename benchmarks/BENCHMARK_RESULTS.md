# Go 结构体映射库性能对比

本文档记录了不同 Go 结构体映射库的性能测试结果。测试环境如下：

- 操作系统：Windows
- 处理器：AMD Ryzen 9 7950X 16-Core Processor
- Go 版本：Go 1.18+

## 测试库

以下是本次测试的映射库：

1. **手动赋值** - 直接通过代码赋值（基准参考）
2. **go-mapster** - 我们自己实现的映射库
3. **jinzhu/copier** - 流行的结构体拷贝库
4. **darccio/mergo** - 结构体和 map 合并库
5. **devfeel/mapper** - 结构体映射库
6. **huandu/go-clone** - 深拷贝库

## 测试场景

测试包含三个场景：

1. **基本结构体映射** - 简单字段映射
2. **嵌套结构体映射** - 包含嵌套结构体和切片的映射
3. **切片映射** - 结构体切片的映射

## 测试结果

### 1. 基本结构体映射

| 库名 | 操作/秒 | 时间/操作 | 内存分配/操作 | 内存分配次数/操作 |
|------|---------|-----------|---------------|-------------------|
| 手动赋值 | 1,000,000,000 | 0.3718 ns/op | 0 B/op | 0 allocs/op |
| go-mapster | 3,061,849 | 390.6 ns/op | 80 B/op | 1 allocs/op |
| darccio/mergo | 3,023,038 | 394.4 ns/op | 328 B/op | 9 allocs/op |
| devfeel/mapper | 1,200,055 | 995.7 ns/op | 392 B/op | 26 allocs/op |
| jinzhu/copier | 682,640 | 1750 ns/op | 536 B/op | 20 allocs/op |
| huandu/go-clone | 12,258,156 | 98.78 ns/op | 160 B/op | 2 allocs/op |

### 2. 嵌套结构体映射

| 库名 | 操作/秒 | 时间/操作 | 内存分配/操作 | 内存分配次数/操作 |
|------|---------|-----------|---------------|-------------------|
| 手动赋值 | 53,618,583 | 21.63 ns/op | 48 B/op | 1 allocs/op |
| huandu/go-clone | 5,565,444 | 215.8 ns/op | 648 B/op | 5 allocs/op |
| go-mapster | 2,346,267 | 512.2 ns/op | 384 B/op | 2 allocs/op |
| jinzhu/copier | 575,138 | 2105 ns/op | 776 B/op | 22 allocs/op |
| darccio/mergo | 1,000,000 | 1144 ns/op | 1248 B/op | 22 allocs/op |
| devfeel/mapper | 942,744 | 1278 ns/op | 984 B/op | 30 allocs/op |

### 3. 切片映射

| 库名 | 操作/秒 | 时间/操作 | 内存分配/操作 | 内存分配次数/操作 |
|------|---------|-----------|---------------|-------------------|
| 手动赋值 | 8,842,794 | 114.9 ns/op | 480 B/op | 1 allocs/op |
| huandu/go-clone | 1,208,168 | 998.2 ns/op | 1488 B/op | 14 allocs/op |
| go-mapster | 454,927 | 2578 ns/op | 552 B/op | 4 allocs/op |
| darccio/mergo | 407,236 | 2950 ns/op | 2800 B/op | 71 allocs/op |
| jinzhu/copier | 345,448 | 3418 ns/op | 1104 B/op | 15 allocs/op |
| devfeel/mapper | 178,783 | 6661 ns/op | 3984 B/op | 118 allocs/op |

## 性能分析

### 基本结构体映射

1. **最快**: 手动赋值 > huandu/go-clone > go-mapster > darccio/mergo > devfeel/mapper > jinzhu/copier
2. **内存效率**: 手动赋值 > go-mapster > huandu/go-clone > darccio/mergo > devfeel/mapper > jinzhu/copier

在基本结构体映射中，手动赋值的性能最好，这是意料之中的。在库中，huandu/go-clone 表现最好，但它主要是克隆而不是映射到不同类型。我们的 go-mapster 在真正的映射库中表现最好，比 darccio/mergo 略快，内存分配也更少。

### 嵌套结构体映射

1. **最快**: 手动赋值 > huandu/go-clone > go-mapster > darccio/mergo > devfeel/mapper > jinzhu/copier
2. **内存效率**: 手动赋值 > go-mapster > huandu/go-clone > jinzhu/copier > devfeel/mapper > darccio/mergo

在嵌套结构体映射中，go-mapster 在真正的映射库中仍然表现最好，内存分配也最少。

### 切片映射

1. **最快**: 手动赋值 > huandu/go-clone > go-mapster > darccio/mergo > jinzhu/copier > devfeel/mapper
2. **内存效率**: 手动赋值 > go-mapster > jinzhu/copier > huandu/go-clone > darccio/mergo > devfeel/mapper

在切片映射中，go-mapster 在内存分配方面表现最好，但在速度上比 huandu/go-clone 慢一些。

## 结论

1. **手动赋值** 始终是最快的方法，但在复杂结构中不实用。
2. **huandu/go-clone** 速度很快，但它主要是用于相同类型的深拷贝，而不是不同类型之间的映射。
3. **go-mapster** 在真正的映射库中表现最好，特别是在基本和嵌套结构体映射方面。
4. **darccio/mergo** 在基本映射中表现良好，但在复杂结构中内存分配较多。
5. **jinzhu/copier** 和 **devfeel/mapper** 在性能上相对较慢，但它们可能提供了其他功能特性。

### 我们的 go-mapster 优势

1. **性能优势**: 在真正的映射库中，go-mapster 在速度和内存效率方面都表现最好。
2. **类型安全**: 使用 Go 1.18+ 泛型提供类型安全。
3. **简洁 API**: 只需一行代码完成映射。
4. **内存效率**: 在所有测试场景中，go-mapster 的内存分配都很低。

### 改进空间

1. **切片映射性能**: 在切片映射方面，go-mapster 的性能可以进一步优化。
2. **减少反射开销**: 可以考虑使用代码生成或缓存来减少反射开销。
