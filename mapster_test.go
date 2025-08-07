package mapster

import (
	"testing"
	"time"
)

// Test structures
type TestUser struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Age       int
	CreatedAt time.Time
}

type TestUserDTO struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	FullName  string
	AgeGroup  string
}

func TestBasicMapping(t *testing.T) {
	user := TestUser{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Age:       30,
		CreatedAt: time.Now(),
	}

	// Test basic mapping without configuration
	dto := Map[TestUserDTO](user)

	if dto.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, dto.ID)
	}
	if dto.FirstName != user.FirstName {
		t.Errorf("Expected FirstName %s, got %s", user.FirstName, dto.FirstName)
	}
	if dto.LastName != user.LastName {
		t.Errorf("Expected LastName %s, got %s", user.LastName, dto.LastName)
	}
	if dto.Email != user.Email {
		t.Errorf("Expected Email %s, got %s", user.Email, dto.Email)
	}

	// FullName and AgeGroup should be empty without configuration
	if dto.FullName != "" {
		t.Errorf("Expected empty FullName, got %s", dto.FullName)
	}
}

func TestMappingWithConfiguration(t *testing.T) {
	// Configure mapping
	Config[TestUser, TestUserDTO]().
		Map("FullName").FromFunc(func(u TestUser) any {
		return u.FirstName + " " + u.LastName
	}).
		Map("AgeGroup").FromFunc(func(u TestUser) any {
		if u.Age < 18 {
			return "Minor"
		} else if u.Age < 65 {
			return "Adult"
		}
		return "Senior"
	}).
		Register()

	user := TestUser{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Age:       30,
		CreatedAt: time.Now(),
	}

	dto := Map[TestUserDTO](user)

	expectedFullName := "John Doe"
	if dto.FullName != expectedFullName {
		t.Errorf("Expected FullName %s, got %s", expectedFullName, dto.FullName)
	}

	expectedAgeGroup := "Adult"
	if dto.AgeGroup != expectedAgeGroup {
		t.Errorf("Expected AgeGroup %s, got %s", expectedAgeGroup, dto.AgeGroup)
	}
}

func TestSliceMapping(t *testing.T) {
	users := []TestUser{
		{ID: 1, FirstName: "John", LastName: "Doe", Age: 30},
		{ID: 2, FirstName: "Jane", LastName: "Smith", Age: 25},
	}

	dtos := make([]TestUserDTO, len(users))
	for i, u := range users {
		dtos[i] = Map[TestUserDTO](u)
	}

	if len(dtos) != len(users) {
		t.Errorf("Expected %d DTOs, got %d", len(users), len(dtos))
	}

	for i, dto := range dtos {
		if dto.ID != users[i].ID {
			t.Errorf("Expected ID %d, got %d", users[i].ID, dto.ID)
		}
		if dto.FirstName != users[i].FirstName {
			t.Errorf("Expected FirstName %s, got %s", users[i].FirstName, dto.FirstName)
		}
	}
}

func TestMapTo(t *testing.T) {
	user := TestUser{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Age:       30,
	}

	var dto TestUserDTO
	MapTo(user, &dto)

	if dto.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, dto.ID)
	}
	if dto.FirstName != user.FirstName {
		t.Errorf("Expected FirstName %s, got %s", user.FirstName, dto.FirstName)
	}
}

func TestNilHandling(t *testing.T) {
	// Test nil source
	dto := Map[TestUserDTO](nil)
	if dto.ID != 0 || dto.FirstName != "" {
		t.Error("Expected zero value for nil source")
	}

	// Test MapTo with nil
	var target TestUserDTO
	target.ID = 999 // Set some value
	MapTo(nil, &target)
	if target.ID != 999 {
		t.Error("MapTo should not modify target when source is nil")
	}
}

// Time conversion test structures
type TimeSource struct {
	ID        int64
	Name      string
	CreatedAt int64 // Unix timestamp
	UpdatedAt int64 // Unix timestamp
}

type TimeTarget struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TimeSourceWithTime struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TimeTargetWithInt64 struct {
	ID        int64
	Name      string
	CreatedAt int64
	UpdatedAt int64
}

func TestTimeConversion(t *testing.T) {
	now := time.Now()
	source := TimeSource{
		ID:        1,
		Name:      "测试用户",
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
	}

	// Test int64 -> time.Time conversion
	target := Map[TimeTarget](source)

	if target.ID != source.ID {
		t.Errorf("ID 转换失败: 期望 %d, 得到 %d", source.ID, target.ID)
	}

	if target.Name != source.Name {
		t.Errorf("Name 转换失败: 期望 %s, 得到 %s", source.Name, target.Name)
	}

	// Verify time conversion
	expectedCreatedAt := time.Unix(source.CreatedAt, 0)
	if !target.CreatedAt.Equal(expectedCreatedAt) {
		t.Errorf("CreatedAt 转换失败: 期望 %v, 得到 %v", expectedCreatedAt, target.CreatedAt)
	}

	expectedUpdatedAt := time.Unix(source.UpdatedAt, 0)
	if !target.UpdatedAt.Equal(expectedUpdatedAt) {
		t.Errorf("UpdatedAt 转换失败: 期望 %v, 得到 %v", expectedUpdatedAt, target.UpdatedAt)
	}
}

func TestTimeConversionReverse(t *testing.T) {
	now := time.Now()
	source := TimeSourceWithTime{
		ID:        2,
		Name:      "测试用户2",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Test time.Time -> int64 conversion
	target := Map[TimeTargetWithInt64](source)

	if target.ID != source.ID {
		t.Errorf("ID 转换失败: 期望 %d, 得到 %d", source.ID, target.ID)
	}

	if target.Name != source.Name {
		t.Errorf("Name 转换失败: 期望 %s, 得到 %s", source.Name, target.Name)
	}

	// Verify time conversion
	expectedCreatedAt := source.CreatedAt.Unix()
	if target.CreatedAt != expectedCreatedAt {
		t.Errorf("CreatedAt 转换失败: 期望 %d, 得到 %d", expectedCreatedAt, target.CreatedAt)
	}

	expectedUpdatedAt := source.UpdatedAt.Unix()
	if target.UpdatedAt != expectedUpdatedAt {
		t.Errorf("UpdatedAt 转换失败: 期望 %d, 得到 %d", expectedUpdatedAt, target.UpdatedAt)
	}
}

func TestTimeConversionConfig(t *testing.T) {
	now := time.Now()
	source := TimeSource{
		ID:        1,
		Name:      "测试用户",
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
	}

	// Test default enabled time conversion
	EnableTimeConversion(true)
	target1 := Map[TimeTarget](source)
	if target1.CreatedAt.IsZero() {
		t.Error("默认启用时间转换应该生效，但得到了零值时间")
	}

	// Test disabled time conversion
	EnableTimeConversion(false)
	target2 := Map[TimeTarget](source)
	if !target2.CreatedAt.IsZero() {
		t.Error("禁用时间转换后应该得到零值时间")
	}

	// Test custom configuration priority
	EnableTimeConversion(true)
	fixedTime := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	Config[TimeSource, TimeTarget]().
		Map("CreatedAt").FromFunc(func(src TimeSource) any {
		return fixedTime
	}).
		Register()

	target3 := Map[TimeTarget](source)
	if !target3.CreatedAt.Equal(fixedTime) {
		t.Errorf("自定义配置应该优先，期望 %v, 得到 %v", fixedTime, target3.CreatedAt)
	}

	// UpdatedAt should use global time conversion
	if target3.UpdatedAt.IsZero() {
		t.Error("UpdatedAt 应该使用全局时间转换")
	}

	// Clean up
	ClearGeneratedMappers()
	globalConfigs = make(map[string]*MappingDefinition)
}
