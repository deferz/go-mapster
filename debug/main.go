package main

import (
	"fmt"
	"reflect"

	"github.com/deferz/go-mapster"
)

type DebugUser struct {
	ID        int64
	FirstName string
	LastName  string
	Age       int
}

type DebugUserDTO struct {
	ID        int64
	FirstName string
	LastName  string
	FullName  string
}

func main() {
	// Check types
	var u DebugUser
	var d DebugUserDTO

	srcType := reflect.TypeOf(u)
	targetType := reflect.TypeOf(d)

	fmt.Println("Source type:", srcType)
	fmt.Println("Target type:", targetType)

	// Configure mapping
	mapster.Config[DebugUser, DebugUserDTO]().
		Map("FullName").FromFunc(func(u DebugUser) interface{} {
		fmt.Println("Custom function called with:", u)
		return u.FirstName + " " + u.LastName
	}).
		Register()

	// Check if config was registered
	config := mapster.GetMappingConfig(srcType, targetType)
	fmt.Println("Config found:", config != nil)
	if config != nil {
		fmt.Printf("Config has %d field mappings\n", len(config.FieldMappings))
		for field, mapping := range config.FieldMappings {
			fmt.Printf("Field %s: Type %d\n", field, mapping.MappingType)
		}
	}

	// Test mapping
	user := DebugUser{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Age:       30,
	}

	fmt.Println("Mapping user:", user)
	dto := mapster.Map[DebugUserDTO](user)
	fmt.Println("Result:", dto)
}
