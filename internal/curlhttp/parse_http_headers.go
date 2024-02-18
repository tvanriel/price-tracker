package curlhttp

import (
	"strings"
)

func ParseHTTPHeaders(source string) map[string][]string {
	headers := make(map[string][]string)
	lines := strings.Split(source, "\n")

	for _, line := range lines {
		// Skip empty lines
		if line == "" {
			continue
		}

		// Skip over the status code line
		if strings.HasPrefix(line, "< HTTP/") {
			continue
		}

		// Skip lines that do not have the incoming sign
		if !strings.HasPrefix(line, "< ") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimPrefix(strings.TrimSpace(parts[0]), "< ")
		value := strings.TrimSpace(parts[1])

		// Append value to existing key or create new entry
		if _, ok := headers[key]; ok {
			headers[key] = append(headers[key], strings.Split(value, ";")...)
		} else {
			headers[key] = strings.Split(value, ";")
		}
	}

	return headers
}
