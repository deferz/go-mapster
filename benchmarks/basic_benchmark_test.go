package benchmarks

import (
	"testing"

	"dario.cat/mergo"
	"github.com/deferz/go-mapster"
	"github.com/devfeel/mapper"
	"github.com/huandu/go-clone"
	"github.com/jinzhu/copier"
)

// 基本结构体定义
type SimpleSource struct {
	ID        int64
	Name      string
	Age       int
	Email     string
	Active    bool
	CreatedAt string
}

type SimpleTarget struct {
	ID        int64
	Name      string
	Age       int
	Email     string
	Active    bool
	CreatedAt string
}

// 初始化测试数据
func getSimpleSource() SimpleSource {
	return SimpleSource{
		ID:        123456789,
		Name:      "张三",
		Age:       30,
		Email:     "zhangsan@example.com",
		Active:    true,
		CreatedAt: "2023-08-09T12:34:56Z",
	}
}

// 手动赋值基准测试
func BenchmarkManualMapping(b *testing.B) {
	src := getSimpleSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dst := SimpleTarget{
			ID:        src.ID,
			Name:      src.Name,
			Age:       src.Age,
			Email:     src.Email,
			Active:    src.Active,
			CreatedAt: src.CreatedAt,
		}
		_ = dst
	}
}

// go-mapster 基准测试
// 注意：映射现在在调用 Map 和 MapTo 时自动注册

func BenchmarkGoMapster(b *testing.B) {
	src := getSimpleSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dst, _ := mapster.Map[SimpleTarget](src)
		_ = dst
	}
}

// jinzhu/copier 基准测试
func BenchmarkJinzhuCopier(b *testing.B) {
	src := getSimpleSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var dst SimpleTarget
		_ = copier.Copy(&dst, &src)
	}
}

// darccio/mergo 基准测试
func BenchmarkDarccioMergo(b *testing.B) {
	src := getSimpleSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var dst SimpleTarget
		_ = mergo.Map(&dst, src)
	}
}

// devfeel/mapper 基准测试
func init() {
	// 初始化 mapper
	mapper.Register(&SimpleSource{})
	mapper.Register(&SimpleTarget{})
}

func BenchmarkDevfeelMapper(b *testing.B) {
	src := getSimpleSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var dst SimpleTarget
		_ = mapper.AutoMapper(&src, &dst)
	}
}

// huandu/go-clone 基准测试
func BenchmarkHuanduClone(b *testing.B) {
	src := getSimpleSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dst := clone.Clone(src).(SimpleSource)
		_ = SimpleTarget{
			ID:        dst.ID,
			Name:      dst.Name,
			Age:       dst.Age,
			Email:     dst.Email,
			Active:    dst.Active,
			CreatedAt: dst.CreatedAt,
		}
	}
}
