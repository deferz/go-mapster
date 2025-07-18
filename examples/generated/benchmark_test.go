package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/deferz/go-mapster"
)

// Benchmark generated vs reflection mapping
func BenchmarkGeneratedMapping(b *testing.B) {
	user := User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Age:       30,
		CreatedAt: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapster.Map[UserDTO](user) // Uses generated mapper
	}
}

func BenchmarkReflectionMapping(b *testing.B) {
	// Temporarily remove generated mapper to force reflection
	mapster.ClearGeneratedMappers()

	user := User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Age:       30,
		CreatedAt: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapster.Map[UserDTO](user) // Uses reflection
	}

	// Re-register generated mapper
	mapster.RegisterGeneratedMapper(mapUserToUserDTO)
}

// Performance comparison function (not a test)
func runPerformanceComparison() {
	// Simple performance test
	user := User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Age:       30,
		CreatedAt: time.Now(),
	}

	// Test generated mapper performance
	start := time.Now()
	for i := 0; i < 100000; i++ {
		_ = mapster.Map[UserDTO](user)
	}
	generatedTime := time.Since(start)

	// Test reflection mapper performance
	mapster.ClearGeneratedMappers()
	start = time.Now()
	for i := 0; i < 100000; i++ {
		_ = mapster.Map[UserDTO](user)
	}
	reflectionTime := time.Since(start)

	// Re-register for next use
	mapster.RegisterGeneratedMapper(mapUserToUserDTO)

	fmt.Printf("Generated mapper (100k ops): %v\n", generatedTime)
	fmt.Printf("Reflection mapper (100k ops): %v\n", reflectionTime)
	fmt.Printf("Speed improvement: %.1fx\n", float64(reflectionTime)/float64(generatedTime))
}
