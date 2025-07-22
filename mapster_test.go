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

	dtos := MapSlice[TestUserDTO](users)

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
