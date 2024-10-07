package cfg

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func EnvInt(env string, def int) int {
	i := def

	if e := os.Getenv(env); e != "" {
		ii, err := strconv.Atoi(e)
		if err == nil {
			i = ii
		}
	}
	return i
}

func EnvInt64(env string, def int64) int64 {
	s := os.Getenv(env)

	if s == "" {
		return def
	}

	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return def
	}

	return v
}

func EnvString(env string, def string) string {
	e := os.Getenv(env)
	if e == "" {
		e = def
	}
	return e
}

func EnvStringArray(env string, def ...string) []string {
	e := os.Getenv(env)
	if e != "" {
		return strings.Fields(e)
	}
	return def
}

func EnvStringFatal(env string) string {
	e := os.Getenv(env)
	if e == "" {
		fmt.Printf("Env %s must be specified\n", env)
		os.Exit(1)
	}
	return e
}

func EnvBool(env string, def bool) bool {
	v := os.Getenv(env)

	if v == "" {
		return def
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}

	return b
}

func EnvDuration(env string, def time.Duration) time.Duration {
	v := os.Getenv(env)

	if v == "" {
		return def
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}

	return d
}

func EnvFloat64(env string, def float64) float64 {
	v := os.Getenv(env)

	if v == "" {
		return def
	}

	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return def
	}

	return f
}

func EnvURLs(env string) []string {
	e := os.Getenv(env)
	if e == "" {
		fmt.Printf("Env %s must be specified\n", env)
		os.Exit(1)
	}

	hosts, err := ParseURLs(e)
	if err != nil {
		fmt.Printf("Env %s could not be parsed: %s\n", env, err)
		os.Exit(1)
	}

	return hosts
}
