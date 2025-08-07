package main

import (
	"fmt"
	"time"

	"github.com/deferz/go-mapster"
)

// Example source structures
type User struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Age       int
	Password  string
	CreatedAt time.Time
}

type Profile struct {
	FirstName string
	LastName  string
	Email     string
	Age       int
}

// Example target structures
type UserDTO struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	FullName  string
	AgeText   string
}

type ProfileDTO struct {
	FirstName string
	LastName  string
	Email     string
	Age       int
	FullName  string
}

func init() {
	// Configure mapping from User to UserDTO
	mapster.Config[User, UserDTO]().
		Map("FullName").FromFunc(func(u User) any {
		return u.FirstName + " " + u.LastName
	}).
		Map("AgeText").FromFunc(func(u User) any {
		return fmt.Sprintf("%d years old", u.Age)
	}).
		Register()

	// Configure mapping from Profile to ProfileDTO
	mapster.Config[Profile, ProfileDTO]().
		Map("FullName").FromFunc(func(p Profile) any {
		return p.FirstName + " " + p.LastName
	}).
		Register()
}

func main() {
	// Test basic mapping
	user := User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Age:       30,
		Password:  "secret123",
		CreatedAt: time.Now(),
	}

	// Map User to UserDTO with custom configuration
	userDTO := mapster.Map[UserDTO](user)
	fmt.Println("User to UserDTO mapping:")
	fmt.Printf("ID: %d\n", userDTO.ID)
	fmt.Printf("FirstName: %s\n", userDTO.FirstName)
	fmt.Printf("LastName: %s\n", userDTO.LastName)
	fmt.Printf("Email: %s\n", userDTO.Email)
	fmt.Printf("FullName: %s\n", userDTO.FullName)
	fmt.Printf("AgeText: %s\n", userDTO.AgeText)
	fmt.Println()

	// Test mapping without configuration (Profile to ProfileDTO)
	profile := Profile{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane.smith@example.com",
		Age:       25,
	}

	profileDTO := mapster.Map[ProfileDTO](profile)
	fmt.Println("Profile to ProfileDTO mapping:")
	fmt.Printf("FirstName: %s\n", profileDTO.FirstName)
	fmt.Printf("LastName: %s\n", profileDTO.LastName)
	fmt.Printf("Email: %s\n", profileDTO.Email)
	fmt.Printf("Age: %d\n", profileDTO.Age)
	fmt.Printf("FullName: %s\n", profileDTO.FullName)
	fmt.Println()

	// Test slice mapping
	users := []User{user, {
		ID:        2,
		FirstName: "Alice",
		LastName:  "Johnson",
		Email:     "alice.johnson@example.com",
		Age:       28,
		Password:  "password456",
	}}

	userDTOs := make([]UserDTO, len(users))
	for i, u := range users {
		userDTOs[i] = mapster.Map[UserDTO](u)
	}
	fmt.Println("Slice mapping:")
	for i, dto := range userDTOs {
		fmt.Printf("User %d: %s (%s)\n", i+1, dto.FullName, dto.AgeText)
	}
	fmt.Println()

	// Test MapTo function
	var existingDTO UserDTO
	mapster.MapTo(user, &existingDTO)
	fmt.Println("MapTo function:")
	fmt.Printf("Mapped to existing object: %s\n", existingDTO.FullName)
}
