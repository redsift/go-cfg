package dcfg

// Version is used to enforce key versioning
type Version uint

// Key is used to address a configuration value
type Key struct {
	Version  Version  // enforce versioned keys
	Type     Type     // type as predictable string
	Elements []string // actual key (first element is the owner)
}

// NewKey creates a Key from version, type, app and additional key parts
func NewKey[T any](v Version, app string, additional ...string) Key {
	return Key{v, TypeOf[T](), append([]string{app}, additional...)}
}
