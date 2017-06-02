package cfg

import (
	"os"
	"testing"
)

func TestParseHostPort(t *testing.T) {
	addr := "tcp://0.0.0.0:49890"

	p, err := ParsePort(addr)
	if err != nil {
		t.Error(err)
	}
	if p != 49890 {
		t.Errorf("expected 49890, got %d", p)
	}

	addr = "tcp://0.0.0.0"
	_, err = ParsePort(addr)
	if err == nil {
		t.Error("expected error when parsing address without port")
	}
}

func TestUpdatePort(t *testing.T) {
	addr := "tcp://0.0.0.0:49890"

	if a := UpdatePort(addr, 123); a != "tcp://0.0.0.0:123" {
		t.Errorf("expected tcp://0.0.0.0:123, got %s", a)
	}

	addr = "tcp://0.0.0.0"
	if a := UpdatePort(addr, 123); a != "tcp://0.0.0.0:123" {
		t.Errorf("expected tcp://0.0.0.0:123, got %s", a)
	}
}

func TestEnvHosts(t *testing.T) {
	addr := "tcp://0.0.0.0:49890-49899"

	os.Setenv("TEST_URLS", addr)
	urls := EnvURLs("TEST_URLS")

	if len(urls) != 10 {
		t.Errorf("Error parsing addr: ", addr)
	}
}
