package redis

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseKeyEvent(t *testing.T) {
	prefixKey := "config"
	testCases := []struct {
		KeyEvent       string
		ExpectedOutput string
	}{
		{
			KeyEvent:       fmt.Sprintf("__keyspace@0__:%s:is-maintenance", prefixKey),
			ExpectedOutput: "is-maintenance",
		},
		{
			KeyEvent:       fmt.Sprintf("__keyspace@0__:%s:test_config", prefixKey),
			ExpectedOutput: "test_config",
		},
		{
			KeyEvent:       fmt.Sprintf("__keyspace@0__:%s:TestConfig", prefixKey),
			ExpectedOutput: "TestConfig",
		},
		{
			KeyEvent:       "__keyspace@0__:aoskdoaskdoaksdokasd",
			ExpectedOutput: "",
		},
	}

	for _, test := range testCases {
		output := parseKeyEvent(prefixKey, test.KeyEvent)
		assert.Equal(t, test.ExpectedOutput, output)
	}
}
