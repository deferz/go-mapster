package mapster

import (
	"testing"
	"time"
)

// Benchmark structures
type BenchUser struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Age       int
	CreatedAt time.Time
	IsActive  bool
}

type BenchUserDTO struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Age       int
	IsActive  bool
	FullName  string
}

func init() {
	// Configure mapping for benchmark
	Config[BenchUser, BenchUserDTO]().
		Map("FullName").FromFunc(func(u BenchUser) any {
		return u.FirstName + " " + u.LastName
	}).
		Register()
}

func BenchmarkBasicMapping(b *testing.B) {
	user := BenchUser{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Age:       30,
		CreatedAt: time.Now(),
		IsActive:  true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Map[BenchUserDTO](user)
	}
}

func BenchmarkSliceMapping(b *testing.B) {
	users := make([]BenchUser, 100)
	for i := 0; i < 100; i++ {
		users[i] = BenchUser{
			ID:        int64(i),
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			Age:       30,
			CreatedAt: time.Now(),
			IsActive:  true,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 使用循环进行批量映射
		dtos := make([]BenchUserDTO, len(users))
		for j, u := range users {
			dtos[j] = Map[BenchUserDTO](u)
		}
	}
}

func BenchmarkMappingWithoutConfiguration(b *testing.B) {
	type SimpleSource struct {
		ID   int64
		Name string
		Age  int
	}

	type SimpleTarget struct {
		ID   int64
		Name string
		Age  int
	}

	source := SimpleSource{
		ID:   1,
		Name: "John",
		Age:  30,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Map[SimpleTarget](source)
	}
}

// 手动映射函数（作为性能基准）
func manualMapBenchUser(src BenchUser) BenchUserDTO {
	return BenchUserDTO{
		ID:        src.ID,
		FirstName: src.FirstName,
		LastName:  src.LastName,
		Email:     src.Email,
		FullName:  src.FirstName + " " + src.LastName,
		Age:       src.Age,
		IsActive:  src.IsActive,
	}
}

// 零反射映射函数
func mapBenchUserToBenchUserDTO(src BenchUser) BenchUserDTO {
	return BenchUserDTO{
		ID:        src.ID,
		FirstName: src.FirstName,
		LastName:  src.LastName,
		Email:     src.Email,
		FullName:  src.FirstName + " " + src.LastName,
		Age:       src.Age,
		IsActive:  src.IsActive,
	}
}

// BenchmarkManualMapping 测试手动映射性能（最快）
func BenchmarkManualMapping(b *testing.B) {
	user := BenchUser{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Age:       30,
		CreatedAt: time.Now(),
		IsActive:  true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manualMapBenchUser(user)
	}
}

// BenchmarkZeroReflectionMapping 测试零反射映射性能
func BenchmarkZeroReflectionMapping(b *testing.B) {
	// 注册零反射映射
	RegisterGeneratedMapper(mapBenchUserToBenchUserDTO)

	user := BenchUser{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Age:       30,
		CreatedAt: time.Now(),
		IsActive:  true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Map[BenchUserDTO](user)
	}
}

// BenchmarkReflectionMapping 测试纯反射映射性能（无配置）
func BenchmarkReflectionMapping(b *testing.B) {
	// 清除注册的映射器，强制使用反射
	ClearGeneratedMappers()

	user := BenchUser{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Age:       30,
		CreatedAt: time.Now(),
		IsActive:  true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Map[BenchUserDTO](user)
	}
}

// BenchmarkPerformanceComparison 综合性能对比
func BenchmarkPerformanceComparison(b *testing.B) {
	user := BenchUser{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Age:       30,
		CreatedAt: time.Now(),
		IsActive:  true,
	}

	// 注册零反射映射
	RegisterGeneratedMapper(mapBenchUserToBenchUserDTO)

	b.Run("1_ManualMapping", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = manualMapBenchUser(user)
		}
	})

	b.Run("2_ZeroReflection", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = Map[BenchUserDTO](user)
		}
	})

	b.Run("3_WithConfig", func(b *testing.B) {
		// 临时清除生成的映射器，使用配置
		ClearGeneratedMappers()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = Map[BenchUserDTO](user)
		}
		// 恢复
		RegisterGeneratedMapper(mapBenchUserToBenchUserDTO)
	})

	b.Run("4_PureReflection", func(b *testing.B) {
		// 使用新类型，无配置无生成器
		type TempUser struct {
			ID    int64
			Name  string
			Email string
		}
		type TempUserDTO struct {
			ID    int64
			Name  string
			Email string
		}

		tempUser := TempUser{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		}

		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = Map[TempUserDTO](tempUser)
		}
	})
}
