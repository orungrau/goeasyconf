package goeasyconf

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// FillConfig populates a struct with values from environment variables.
// The struct fields should be tagged with `env` for the variable name and `required:"true"` if necessary.
// cfg should be a pointer to a struct; otherwise, an error is returned.
func FillConfig(cfg interface{}) error {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("cfg must be a pointer to a struct")
	}
	return populateStruct(v.Elem())
}

// populateStruct recursively populates fields in a struct with environment variables based on struct tags.
func populateStruct(v reflect.Value) error {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		envKey := fieldType.Tag.Get("env")
		required := fieldType.Tag.Get("required") == "true"

		// Recursively populate nested structs
		if field.Kind() == reflect.Struct {
			if err := populateStruct(field); err != nil {
				return err
			}
			continue
		}

		envValue := os.Getenv(envKey)

		// Check if required environment variable is missing
		if required && envValue == "" {
			return errors.New("missing required environment variable: " + envKey)
		}

		// Set the field's value if the environment variable is present
		if envValue != "" {
			if err := setFieldValue(field, envValue); err != nil {
				return err
			}
		}
	}
	return nil
}

// setFieldValue converts a string from an environment variable into the appropriate type and sets it to the field.
func setFieldValue(field reflect.Value, envValue string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(envValue)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(envValue, 10, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetInt(intVal)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(envValue, 10, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetUint(uintVal)

	case reflect.Bool:
		boolVal, err := strconv.ParseBool(envValue)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)

	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(envValue, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetFloat(floatVal)

	case reflect.Slice:
		values := strings.Split(envValue, ",")
		slice := reflect.MakeSlice(field.Type(), len(values), len(values))

		for i, val := range values {
			elem := slice.Index(i)
			if err := setFieldValue(elem, val); err != nil {
				return err
			}
		}
		field.Set(slice)

	default:
		return errors.New("unsupported field type: " + field.Kind().String())
	}
	return nil
}
