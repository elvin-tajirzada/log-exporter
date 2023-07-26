package utils

import (
	"fmt"
	"regexp"
)

func ExtractJSON(value string) (string, error) {
	re := regexp.MustCompile(`{.*}`)

	match := re.FindString(value)

	if match == "" {
		return "", fmt.Errorf("unable to find JSON data in log line. line: %s", value)
	}

	return match, nil
}
