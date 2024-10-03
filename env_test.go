package cfg_test

import (
	"os"
	"testing"
	"time"

	"github.com/redsift/go-cfg"
)

func TestEnvInt(t *testing.T) {
	tests := []struct {
		name     string
		env      string
		def      int
		envValue string
		want     int
	}{
		{"Default", "TEST_ENV_INT", 42, "", 42},
		{"ValidValue", "TEST_ENV_INT", 0, "10", 10},
		{"InvalidValue", "TEST_ENV_INT", 42, "not_a_number", 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.env, tt.envValue)
			defer os.Unsetenv(tt.env)

			if got := cfg.EnvInt(tt.env, tt.def); got != tt.want {
				t.Errorf("EnvInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvString(t *testing.T) {
	tests := []struct {
		name     string
		env      string
		def      string
		envValue string
		want     string
	}{
		{"Default", "TEST_ENV_STRING", "default", "", "default"},
		{"EnvValue", "TEST_ENV_STRING", "default", "custom", "custom"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.env, tt.envValue)
			defer os.Unsetenv(tt.env)

			if got := cfg.EnvString(tt.env, tt.def); got != tt.want {
				t.Errorf("EnvString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvStringArray(t *testing.T) {
	tests := []struct {
		name     string
		env      string
		def      []string
		envValue string
		want     []string
	}{
		{"Default", "TEST_ENV_STRING_ARRAY", []string{"a", "b"}, "", []string{"a", "b"}},
		{"EnvValue", "TEST_ENV_STRING_ARRAY", []string{}, "x y z", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.env, tt.envValue)
			defer os.Unsetenv(tt.env)

			got := cfg.EnvStringArray(tt.env, tt.def...)
			if len(got) != len(tt.want) {
				t.Errorf("EnvStringArray() = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("EnvStringArray()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestEnvBool(t *testing.T) {
	tests := []struct {
		name     string
		env      string
		def      bool
		envValue string
		want     bool
	}{
		{"Default", "TEST_ENV_BOOL", true, "", true},
		{"TrueValue", "TEST_ENV_BOOL", false, "true", true},
		{"FalseValue", "TEST_ENV_BOOL", true, "false", false},
		{"InvalidValue", "TEST_ENV_BOOL", true, "not_a_bool", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.env, tt.envValue)
			defer os.Unsetenv(tt.env)

			if got := cfg.EnvBool(tt.env, tt.def); got != tt.want {
				t.Errorf("EnvBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvDuration(t *testing.T) {
	tests := []struct {
		name     string
		env      string
		def      time.Duration
		envValue string
		want     time.Duration
	}{
		{"Default", "TEST_ENV_DURATION", 5 * time.Second, "", 5 * time.Second},
		{"ValidValue", "TEST_ENV_DURATION", time.Second, "10s", 10 * time.Second},
		{"InvalidValue", "TEST_ENV_DURATION", time.Minute, "not_a_duration", time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.env, tt.envValue)
			defer os.Unsetenv(tt.env)

			if got := cfg.EnvDuration(tt.env, tt.def); got != tt.want {
				t.Errorf("EnvDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvFloat64(t *testing.T) {
	tests := []struct {
		name     string
		env      string
		def      float64
		envValue string
		want     float64
	}{
		{"Default", "TEST_ENV_FLOAT", 3.14, "", 3.14},
		{"ValidValue", "TEST_ENV_FLOAT", 0, "2.718", 2.718},
		{"InvalidValue", "TEST_ENV_FLOAT", 1.23, "not_a_float", 1.23},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.env, tt.envValue)
			defer os.Unsetenv(tt.env)

			if got := cfg.EnvFloat64(tt.env, tt.def); got != tt.want {
				t.Errorf("EnvFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}
