package benchmarks

import (
	"testing"

	"dario.cat/mergo"
	"github.com/deferz/go-mapster"
	"github.com/devfeel/mapper"
	"github.com/huandu/go-clone"
	"github.com/jinzhu/copier"
)

// 集合结构体定义
type User struct {
	ID    int
	Name  string
	Email string
	Age   int
}

type UserDTO struct {
	ID    int
	Name  string
	Email string
	Age   int
}

// 初始化测试数据
func getUsers(count int) []User {
	users := make([]User, count)
	for i := 0; i < count; i++ {
		users[i] = User{
			ID:    i + 1,
			Name:  "用户" + string(rune(i+65)),
			Email: "user" + string(rune(i+65)) + "@example.com",
			Age:   20 + i%10,
		}
	}
	return users
}

// 手动赋值基准测试
func BenchmarkSliceManualMapping(b *testing.B) {
	src := getUsers(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dst := make([]UserDTO, len(src))
		for j, user := range src {
			dst[j] = UserDTO{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
				Age:   user.Age,
			}
		}
		_ = dst
	}
}

// go-mapster 基准测试
func BenchmarkSliceGoMapster(b *testing.B) {
	src := getUsers(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dst, _ := mapster.Map[[]UserDTO](src)
		_ = dst
	}
}

// jinzhu/copier 基准测试
func BenchmarkSliceJinzhuCopier(b *testing.B) {
	src := getUsers(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var dst []UserDTO
		_ = copier.Copy(&dst, &src)
	}
}

// darccio/mergo 基准测试 (mergo 不直接支持切片映射，需要手动循环)
func BenchmarkSliceDarccioMergo(b *testing.B) {
	src := getUsers(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dst := make([]UserDTO, len(src))
		for j, user := range src {
			_ = mergo.Map(&dst[j], user)
		}
		_ = dst
	}
}

// devfeel/mapper 基准测试
func init() {
	// 初始化 mapper
	mapper.Register(&User{})
	mapper.Register(&UserDTO{})
}

func BenchmarkSliceDevfeelMapper(b *testing.B) {
	src := getUsers(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var dst []UserDTO
		_ = mapper.MapperSlice(src, &dst)
	}
}

// huandu/go-clone 基准测试
func BenchmarkSliceHuanduClone(b *testing.B) {
	src := getUsers(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cloned := clone.Clone(src).([]User)
		dst := make([]UserDTO, len(cloned))
		for j, user := range cloned {
			dst[j] = UserDTO{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
				Age:   user.Age,
			}
		}
		_ = dst
	}
}
