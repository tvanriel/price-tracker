package curlhttp

import (
	"errors"
	"regexp"
	"strings"
)

type HTTPStatusCodeParser func(source string) (code string, err error)

func CurlHttpStatusCodeParser(source string) (string, error) {
	line, found := findHttpLine(strings.Split(source, "\n"))
	if !found {
		return "", errors.New("cannot find http status response header")
	}
	code, found := getCodeFromHttpLine(line)
	if !found {
		return "", errors.New("cannot find status code in status header")
	}
	return code, nil
}

func findHttpLine(lines []string) (string, bool) {
	for i := range lines {
		matched, err := regexp.MatchString(`< HTTP/[\d.]+ [\d]+[a-zA-Z ]*`, lines[i])
		if err != nil {
			return "", false
		}
		if matched {
			return lines[i], true
		}
	}
	return "", false
}

func getCodeFromHttpLine(line string) (string, bool) {
	// Remove everything before the code.
	_, code, found := strings.Cut(strings.Replace(line, "< HTTP/", "", 1), " ")
	if !found {
		return code, found
	}
	// Remove everything after the code
	code, _, found = strings.Cut(code, " ")
	return code, found
}
