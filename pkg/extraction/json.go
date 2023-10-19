package extraction

import (
	"fmt"
	"regexp"
)

const jsonExpression = `{.*}`

func JSON(value string) (string, error) {
	re := regexp.MustCompile(jsonExpression)

	match := re.FindString(value)

	if match == "" {
		return "", fmt.Errorf("unable to find JSON data in log line. line: %s", value)
	}

	return match, nil
}
