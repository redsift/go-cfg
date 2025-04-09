package mongodb

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const manglePrefix = "mangled__"

var (
	demangleKeyRE = regexp.MustCompile("_[a-fA-F0-9][a-fA-F0-9]")
	mangleKeyRE   = regexp.MustCompile("[^a-zA-Z0-9]")
)

// demangleKey decodes a key from a mongo-safe string to the source string.
func (m *Map) demangleKey(key string) (result string, demangled bool, err error) {
	result, demangled = strings.CutPrefix(key, manglePrefix)

	if !demangled {
		return key, demangled, nil
	}

	var errs []error

	result = demangleKeyRE.ReplaceAllStringFunc(result, func(in string) string {
		v, err := strconv.ParseUint(in[1:], 16, 8)
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid key %q, contains invalid escape sequence %q", key, in))
			return in
		}
		return string(byte(v))
	})

	return result, true, errors.Join(errs...)
}

// mangleKey encodes a string to be a mongo-safe property name.
func (m *Map) mangleKey(key string) string {
	n := 0
	mangled := mangleKeyRE.ReplaceAllStringFunc(key, func(in string) string {
		n++
		return fmt.Sprintf("_%02x", in)
	})
	if key == "" || n > 0 {
		return manglePrefix + mangled
	}
	return key
}
