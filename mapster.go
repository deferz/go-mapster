package mapster

import (
	"fmt"
	"reflect"

	"github.com/deferz/go-mapster/internal/mapper"
)

// Map maps the source object to the target type and returns a new instance.
// This function uses generics to ensure type safety, checking type matches at compile time.
//
// Parameters:
//   - src: Source object, can be any type
//
// Returns:
//   - T: New instance of the target type
//   - error: Returns an error if mapping fails
//
// Example:
//
//	userDTO, err := mapster.Map[UserDTO](user)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Supported mapping types:
//   - Struct to struct
//   - Conversion between basic types
//   - Slices and arrays
//   - Pointer types
//   - Nested structs
func Map[T any](src any) (T, error) {
	var result T

	if src == nil {
		return result, fmt.Errorf("source cannot be nil")
	}

	// Create pointer to target type
	resultPtr := &result

	// Execute mapping
	if err := mapper.MapValue(reflect.ValueOf(src), reflect.ValueOf(resultPtr).Elem()); err != nil {
		return result, fmt.Errorf("mapping failed: %w", err)
	}

	return result, nil
}

// MapTo maps the source object to an existing target object.
// This function modifies the provided target object instead of creating a new instance.
//
// Parameters:
//   - src: Source object, can be any type
//   - dst: Pointer to target object, mapping result will be written to this object
//
// Returns:
//   - error: Returns an error if mapping fails
//
// Example:
//
//	var userDTO UserDTO
//	err := mapster.MapTo(user, &userDTO)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Notes:
//   - dst must be a pointer type
//   - Existing field values may be overwritten
//   - If source object doesn't have a corresponding field, target object's field will retain its original value
func MapTo[T any](src any, dst *T) error {
	if src == nil {
		return fmt.Errorf("source cannot be nil")
	}

	if dst == nil {
		return fmt.Errorf("destination cannot be nil")
	}

	// Execute mapping
	if err := mapper.MapValue(reflect.ValueOf(src), reflect.ValueOf(dst).Elem()); err != nil {
		return fmt.Errorf("mapping failed: %w", err)
	}

	return nil
}
