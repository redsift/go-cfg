package cfg

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// parsePort returns the port parsed from addr or an error if there is no port
// in addr.
func ParsePort(addr string) (int, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return 0, err
	}
	_, ps, err := net.SplitHostPort(u.Host)
	if err != nil {
		return 0, err
	}
	p, err := strconv.ParseInt(ps, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(p), nil
}

// update port returns an address with updated port.
func UpdatePort(addr string, port int) string {
	u, err := url.Parse(addr)
	if err != nil {
		return ""
	}

	h, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		// Assuming address doesn't contain port.
		h = u.Host
	}
	u.Host = net.JoinHostPort(h, strconv.FormatInt(int64(port), 10))
	return u.String()
}

func ParseURLs(addrs string) ([]string, error) {
	var hosts []string
	// Parse tcp:0.0.0.0:port1,tcp:0.0.0.0:port2,...
	hosts = strings.Split(addrs, ",")

	if len(hosts) == 1 {
		// Parse tcp://0.0.0.0:startPort OR tcp://0.0.0.0:startPort-endPort

		splitColon := strings.Split(addrs, ":")
		if len(splitColon) >= 3 {
			rangeStr := splitColon[2]
			ranges := strings.Split(rangeStr, "-")
			if len(ranges) != 2 {
				ranges = []string{rangeStr, rangeStr}
			}

			startPort, err := strconv.Atoi(ranges[0])
			if err != nil {
				return nil, err
			}

			stopPort, err := strconv.Atoi(ranges[1])
			if err != nil {
				return nil, err
			}

			if stopPort < startPort {
				return nil, fmt.Errorf("Invalid port range")
			}

			hosts = nil
			for i := startPort; i <= stopPort; i++ {
				addr := splitColon[0] + ":" + splitColon[1] + ":" + strconv.Itoa(i)
				hosts = append(hosts, addr)
			}
		}
	}

	if len(hosts) == 0 {
		return nil, fmt.Errorf("You must specify atleast 1 url")
	}

	return hosts, nil
}
