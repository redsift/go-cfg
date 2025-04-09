package dcfg

import (
	"fmt"
	"net/url"
	"os"
)

const (
	EnvDCFGBackendURI = "DCFG_BACKEND_URI"
)

func NewFromEnv() (Backend, error) {
	uri, found := os.LookupEnv(EnvDCFGBackendURI)
	if !found {
		return nil, fmt.Errorf("environment variable %q is not set", EnvDCFGBackendURI)
	}

	parsed, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	factory, ok := impls[parsed.Scheme]
	if !ok {
		return nil, fmt.Errorf("Backend for %q not configured", parsed.Scheme)
	}

	return factory(parsed)
}
