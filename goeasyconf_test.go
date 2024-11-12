package goeasyconf

import (
	"os"
	"reflect"
	"testing"
)

type NestedConfig struct {
	NestedInt   int    `env:"NESTED_INT"`
	NestedBool  bool   `env:"NESTED_BOOL"`
	NestedSlice []int  `env:"NESTED_SLICE"`
	NestedStr   string `env:"NESTED_STR"`
}

type Config struct {
	StringValue  string   `env:"STRING_VALUE" required:"true"`
	IntValue     int      `env:"INT_VALUE"`
	UintValue    uint     `env:"UINT_VALUE"`
	BoolValue    bool     `env:"BOOL_VALUE"`
	FloatValue   float64  `env:"FLOAT_VALUE"`
	StringSlice  []string `env:"STRING_SLICE"`
	IntSlice     []int    `env:"INT_SLICE"`
	NestedConfig NestedConfig
}

// TestFillConfig validates that FillConfig correctly populates fields from environment variables.
func TestFillConfig(t *testing.T) {
	// Set up environment variables
	os.Setenv("STRING_VALUE", "test_string")
	os.Setenv("INT_VALUE", "42")
	os.Setenv("UINT_VALUE", "100")
	os.Setenv("BOOL_VALUE", "true")
	os.Setenv("FLOAT_VALUE", "3.14")
	os.Setenv("STRING_SLICE", "apple,banana,orange")
	os.Setenv("INT_SLICE", "1,2,3")
	os.Setenv("NESTED_INT", "99")
	os.Setenv("NESTED_BOOL", "false")
	os.Setenv("NESTED_SLICE", "10,20,30")
	os.Setenv("NESTED_STR", "nested_test")

	defer func() {
		// Clean up environment variables
		os.Unsetenv("STRING_VALUE")
		os.Unsetenv("INT_VALUE")
		os.Unsetenv("UINT_VALUE")
		os.Unsetenv("BOOL_VALUE")
		os.Unsetenv("FLOAT_VALUE")
		os.Unsetenv("STRING_SLICE")
		os.Unsetenv("INT_SLICE")
		os.Unsetenv("NESTED_INT")
		os.Unsetenv("NESTED_BOOL")
		os.Unsetenv("NESTED_SLICE")
		os.Unsetenv("NESTED_STR")
	}()

	// Create a config instance and call FillConfig
	cfg := Config{}
	err := FillConfig(&cfg)
	if err != nil {
		t.Fatalf("FillConfig failed: %v", err)
	}

	// Expected config for comparison
	expectedCfg := Config{
		StringValue: "test_string",
		IntValue:    42,
		UintValue:   100,
		BoolValue:   true,
		FloatValue:  3.14,
		StringSlice: []string{"apple", "banana", "orange"},
		IntSlice:    []int{1, 2, 3},
		NestedConfig: NestedConfig{
			NestedInt:   99,
			NestedBool:  false,
			NestedSlice: []int{10, 20, 30},
			NestedStr:   "nested_test",
		},
	}

	// Compare the expected and actual configs
	if !reflect.DeepEqual(cfg, expectedCfg) {
		t.Errorf("Config does not match expected values.\nGot: %#v\nExpected: %#v", cfg, expectedCfg)
	}
}
