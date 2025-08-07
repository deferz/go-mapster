package main

import (
	"fmt"
	"time"

	"github.com/deferz/go-mapster"
)

// 示例源结构
type User struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Age       int
	Password  string
	CreatedAt time.Time
}

// 示例目标结构
type UserDTO struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	FullName  string
	AgeText   string
	IsAdult   bool
}

// 手动编写的映射函数 - 零反射，最高性能
func mapUserToUserDTO(src User) UserDTO {
	return UserDTO{
		ID:        src.ID,
		FirstName: src.FirstName,
		LastName:  src.LastName,
		Email:     src.Email,
		FullName:  src.FirstName + " " + src.LastName,
		AgeText:   fmt.Sprintf("%d 岁", src.Age),
		IsAdult:   src.Age >= 18,
	}
}

// 另一个手动映射函数示例
type Profile struct {
	FirstName string
	LastName  string
	Email     string
	Age       int
}

type ProfileDTO struct {
	FirstName string
	LastName  string
	Email     string
	Age       int
	FullName  string
}

func mapProfileToProfileDTO(src Profile) ProfileDTO {
	return ProfileDTO{
		FirstName: src.FirstName,
		LastName:  src.LastName,
		Email:     src.Email,
		Age:       src.Age,
		FullName:  src.FirstName + " " + src.LastName,
	}
}

func init() {
	// 注册手动编写的映射函数
	mapster.RegisterGeneratedMapper(mapUserToUserDTO)
	mapster.RegisterGeneratedMapper(mapProfileToProfileDTO)
}

func main() {
	fmt.Println("=== 手动零反射映射示例 ===")

	// 测试用户映射
	user := User{
		ID:        1,
		FirstName: "张",
		LastName:  "三",
		Email:     "zhangsan@example.com",
		Age:       25,
		Password:  "secret123",
		CreatedAt: time.Now(),
	}

	// 使用手动映射函数（零反射，最快）
	userDTO := mapster.Map[UserDTO](user)
	fmt.Printf("用户映射结果: %+v\n", userDTO)

	// 测试配置文件映射
	profile := Profile{
		FirstName: "李",
		LastName:  "四",
		Email:     "lisi@example.com",
		Age:       30,
	}

	// 使用手动映射函数
	profileDTO := mapster.Map[ProfileDTO](profile)
	fmt.Printf("配置文件映射结果: %+v\n", profileDTO)

	fmt.Println("\n=== 性能对比 ===")

	// 简单性能测试
	iterations := 100000

	// 测试手动映射性能
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = mapster.Map[UserDTO](user)
	}
	manualTime := time.Since(start)

	// 测试配置映射性能（临时清除手动映射器）
	mapster.ClearGeneratedMappers()
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = mapster.Map[UserDTO](user)
	}
	configTime := time.Since(start)

	// 恢复手动映射器
	mapster.RegisterGeneratedMapper(mapUserToUserDTO)

	fmt.Printf("手动映射 (%dk 次): %v\n", iterations/1000, manualTime)
	fmt.Printf("配置映射 (%dk 次): %v\n", iterations/1000, configTime)
	fmt.Printf("性能提升: %.1fx\n", float64(configTime)/float64(manualTime))
}
