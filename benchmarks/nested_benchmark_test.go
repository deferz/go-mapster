package benchmarks

import (
	"testing"

	"dario.cat/mergo"
	"github.com/deferz/go-mapster"
	"github.com/devfeel/mapper"
	"github.com/huandu/go-clone"
	"github.com/jinzhu/copier"
)

// 嵌套结构体定义
type Address struct {
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
}

type ContactInfo struct {
	Phone   string
	Email   string
	Website string
}

type NestedSource struct {
	ID          int64
	Name        string
	Age         int
	Address     Address
	ContactInfo ContactInfo
	Tags        []string
	Active      bool
}

type NestedTarget struct {
	ID          int64
	Name        string
	Age         int
	Address     Address
	ContactInfo ContactInfo
	Tags        []string
	Active      bool
}

// 初始化测试数据
func getNestedSource() NestedSource {
	return NestedSource{
		ID:     123456789,
		Name:   "张三",
		Age:    30,
		Active: true,
		Address: Address{
			Street:     "人民路123号",
			City:       "上海",
			State:      "上海",
			PostalCode: "200001",
			Country:    "中国",
		},
		ContactInfo: ContactInfo{
			Phone:   "13800138000",
			Email:   "zhangsan@example.com",
			Website: "https://zhangsan.example.com",
		},
		Tags: []string{"用户", "VIP", "活跃"},
	}
}

// 手动赋值基准测试
func BenchmarkNestedManualMapping(b *testing.B) {
	src := getNestedSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dst := NestedTarget{
			ID:     src.ID,
			Name:   src.Name,
			Age:    src.Age,
			Active: src.Active,
			Address: Address{
				Street:     src.Address.Street,
				City:       src.Address.City,
				State:      src.Address.State,
				PostalCode: src.Address.PostalCode,
				Country:    src.Address.Country,
			},
			ContactInfo: ContactInfo{
				Phone:   src.ContactInfo.Phone,
				Email:   src.ContactInfo.Email,
				Website: src.ContactInfo.Website,
			},
			Tags: append([]string{}, src.Tags...),
		}
		_ = dst
	}
}

// go-mapster 基准测试
func BenchmarkNestedGoMapster(b *testing.B) {
	src := getNestedSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dst, _ := mapster.Map[NestedTarget](src)
		_ = dst
	}
}

// jinzhu/copier 基准测试
func BenchmarkNestedJinzhuCopier(b *testing.B) {
	src := getNestedSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var dst NestedTarget
		_ = copier.Copy(&dst, &src)
	}
}

// darccio/mergo 基准测试
func BenchmarkNestedDarccioMergo(b *testing.B) {
	src := getNestedSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var dst NestedTarget
		_ = mergo.Map(&dst, src)
	}
}

// devfeel/mapper 基准测试
func init() {
	// 初始化 mapper
	mapper.Register(&NestedSource{})
	mapper.Register(&NestedTarget{})
	mapper.Register(&Address{})
	mapper.Register(&ContactInfo{})
}

func BenchmarkNestedDevfeelMapper(b *testing.B) {
	src := getNestedSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var dst NestedTarget
		_ = mapper.AutoMapper(&src, &dst)
	}
}

// huandu/go-clone 基准测试
func BenchmarkNestedHuanduClone(b *testing.B) {
	src := getNestedSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cloned := clone.Clone(src).(NestedSource)
		dst := NestedTarget{
			ID:          cloned.ID,
			Name:        cloned.Name,
			Age:         cloned.Age,
			Active:      cloned.Active,
			Address:     cloned.Address,
			ContactInfo: cloned.ContactInfo,
			Tags:        cloned.Tags,
		}
		_ = dst
	}
}
