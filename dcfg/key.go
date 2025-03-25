package dcfg

// Version is used to enforce key versioning
type Version uint

// Key is used to address a configuration value
type Key struct {
	Version  Version
	Elements []string
}

// NewKey creates a Key from version, app and additional key parts
func NewKey(v Version, app string, additional ...string) Key {
	return Key{v, append([]string{app}, additional...)}
}
