package redis

import (
	"fmt"
	"regexp"
)

func parseKeyEvent(prefixKey, event string) (affectedKey string) {
	eventRegexString := fmt.Sprintf(`(?m)^__keyspace@0__:%s:(?P<key>[^\s]+)$`, prefixKey)
	eventRegex := regexp.MustCompile(eventRegexString)

	match := eventRegex.FindStringSubmatch(event)
	result := map[string]string{}

	if len(match) != len(eventRegex.SubexpNames()) {
		return ""
	}

	for i, name := range eventRegex.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}

	return result["key"]
}
