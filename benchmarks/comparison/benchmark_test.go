package comparison

import (
	"testing"
	"time"

	"github.com/deferz/go-mapster"
	"github.com/devfeel/mapper"
	"github.com/jinzhu/copier"
	"github.com/mitchellh/mapstructure"
)

// 测试用的源结构体
type User struct {
	ID          int64
	FirstName   string
	LastName    string
	Email       string
	Age         int
	Phone       string
	Address     string
	City        string
	Country     string
	PostalCode  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	IsActive    bool
	Score       float64
	Tags        []string
	Preferences map[string]string
}

// 目标结构体
type UserDTO struct {
	ID         int64
	FirstName  string
	LastName   string
	Email      string
	FullName   string // 需要组合 FirstName + LastName
	Age        int
	Phone      string
	Address    string
	City       string
	Country    string
	PostalCode string
	IsActive   bool
	Score      float64
}

// 简单结构体（用于基础测试）
type SimpleUser struct {
	ID    int64
	Name  string
	Email string
	Age   int
}

type SimpleUserDTO struct {
	ID    int64
	Name  string
	Email string
	Age   int
}

// 初始化
func init() {
	// 配置 mapster
	mapster.Config[User, UserDTO]().
		Map("FullName").FromFunc(func(u User) any {
		return u.FirstName + " " + u.LastName
	}).
		Register()
}

// 测试数据
var testUser = User{
	ID:         123,
	FirstName:  "John",
	LastName:   "Doe",
	Email:      "john.doe@example.com",
	Age:        30,
	Phone:      "123-456-7890",
	Address:    "123 Main St",
	City:       "New York",
	Country:    "USA",
	PostalCode: "10001",
	CreatedAt:  time.Now(),
	UpdatedAt:  time.Now(),
	IsActive:   true,
	Score:      98.5,
	Tags:       []string{"vip", "premium"},
	Preferences: map[string]string{
		"theme":    "dark",
		"language": "en",
	},
}

var simpleUser = SimpleUser{
	ID:    1,
	Name:  "John Doe",
	Email: "john@example.com",
	Age:   30,
}

// BenchmarkSimpleMapping 测试简单对象映射
func BenchmarkSimpleMapping(b *testing.B) {
	b.Run("Mapster", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = mapster.Map[SimpleUserDTO](simpleUser)
		}
	})

	b.Run("Copier", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var dto SimpleUserDTO
			_ = copier.Copy(&dto, &simpleUser)
		}
	})

	b.Run("DevfeelMapper", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dto := &SimpleUserDTO{}
			_ = mapper.AutoMapper(&simpleUser, dto)
		}
	})

	b.Run("Mapstructure", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var dto SimpleUserDTO
			_ = mapstructure.Decode(simpleUser, &dto)
		}
	})

	b.Run("Manual", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = SimpleUserDTO{
				ID:    simpleUser.ID,
				Name:  simpleUser.Name,
				Email: simpleUser.Email,
				Age:   simpleUser.Age,
			}
		}
	})
}

// BenchmarkComplexMapping 测试复杂对象映射
func BenchmarkComplexMapping(b *testing.B) {
	b.Run("Mapster", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = mapster.Map[UserDTO](testUser)
		}
	})

	b.Run("Copier", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var dto UserDTO
			_ = copier.Copy(&dto, &testUser)
			// Copier 不支持自定义字段映射，需要手动处理
			dto.FullName = testUser.FirstName + " " + testUser.LastName
		}
	})

	b.Run("DevfeelMapper", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dto := &UserDTO{}
			_ = mapper.AutoMapper(&testUser, dto)
			// 手动处理 FullName
			dto.FullName = testUser.FirstName + " " + testUser.LastName
		}
	})

	b.Run("Manual", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = UserDTO{
				ID:         testUser.ID,
				FirstName:  testUser.FirstName,
				LastName:   testUser.LastName,
				Email:      testUser.Email,
				FullName:   testUser.FirstName + " " + testUser.LastName,
				Age:        testUser.Age,
				Phone:      testUser.Phone,
				Address:    testUser.Address,
				City:       testUser.City,
				Country:    testUser.Country,
				PostalCode: testUser.PostalCode,
				IsActive:   testUser.IsActive,
				Score:      testUser.Score,
			}
		}
	})
}

// BenchmarkSliceMapping 测试切片映射
func BenchmarkSliceMapping(b *testing.B) {
	// 准备测试数据
	users := make([]User, 100)
	for i := range users {
		users[i] = testUser
		users[i].ID = int64(i)
	}

	b.Run("Mapster", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dtos := make([]UserDTO, len(users))
			for j, u := range users {
				dtos[j] = mapster.Map[UserDTO](u)
			}
		}
	})

	b.Run("Copier", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var dtos []UserDTO
			_ = copier.Copy(&dtos, &users)
			// 处理 FullName
			for j := range dtos {
				dtos[j].FullName = users[j].FirstName + " " + users[j].LastName
			}
		}
	})

	b.Run("Manual", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dtos := make([]UserDTO, len(users))
			for j, u := range users {
				dtos[j] = UserDTO{
					ID:         u.ID,
					FirstName:  u.FirstName,
					LastName:   u.LastName,
					Email:      u.Email,
					FullName:   u.FirstName + " " + u.LastName,
					Age:        u.Age,
					Phone:      u.Phone,
					Address:    u.Address,
					City:       u.City,
					Country:    u.Country,
					PostalCode: u.PostalCode,
					IsActive:   u.IsActive,
					Score:      u.Score,
				}
			}
		}
	})
}

// BenchmarkNestedMapping 测试嵌套结构映射
func BenchmarkNestedMapping(b *testing.B) {
	type Address struct {
		Street     string
		City       string
		Country    string
		PostalCode string
	}

	type UserWithAddress struct {
		ID        int64
		Name      string
		Email     string
		Address   Address
		CreatedAt time.Time
	}

	type UserWithAddressDTO struct {
		ID      int64
		Name    string
		Email   string
		Address Address
	}

	testUserWithAddr := UserWithAddress{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
		Address: Address{
			Street:     "123 Main St",
			City:       "New York",
			Country:    "USA",
			PostalCode: "10001",
		},
		CreatedAt: time.Now(),
	}

	b.Run("Mapster", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = mapster.Map[UserWithAddressDTO](testUserWithAddr)
		}
	})

	b.Run("Copier", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var dto UserWithAddressDTO
			_ = copier.Copy(&dto, &testUserWithAddr)
		}
	})

	b.Run("Manual", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = UserWithAddressDTO{
				ID:      testUserWithAddr.ID,
				Name:    testUserWithAddr.Name,
				Email:   testUserWithAddr.Email,
				Address: testUserWithAddr.Address,
			}
		}
	})
}
