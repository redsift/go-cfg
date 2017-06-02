package cfg

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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
