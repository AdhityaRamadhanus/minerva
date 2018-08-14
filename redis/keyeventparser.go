package redis

import (
	"regexp"

	"github.com/AdhityaRamadhanus/minerva"
)

func parseKeyEvent(event, payload string) minerva.KeyEvent {
	eventRegex := regexp.MustCompile(`(?m)^__keyspace@0__:config:(?P<key>[^\s]+)$`)

	match := eventRegex.FindStringSubmatch(event)
	result := map[string]string{}
	for i, name := range eventRegex.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}

	return minerva.KeyEvent{
		Type:        payload,
		AffectedKey: result["key"],
	}
}
