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
		Map("FullName").FromFunc(func(u BenchUser) interface{} {
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
		_ = MapSlice[BenchUserDTO](users)
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
