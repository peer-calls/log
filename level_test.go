package log_test

import (
	"testing"

	"github.com/peer-calls/log"
	"github.com/stretchr/testify/assert"
)

func TestLevel_String(t *testing.T) {
	t.Parallel()

	type testCase struct {
		level      log.Level
		wantString string
	}

	testCases := []testCase{
		{log.LevelError, "error"},
		{log.LevelWarn, "warn"},
		{log.LevelInfo, "info"},
		{log.LevelDebug, "debug"},
		{log.LevelTrace, "trace"},
		{log.LevelDisabled, "disabled"},
		{log.LevelUnknown, "Unknown(-1)"},
		{log.Level(-3), "Unknown(-3)"},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.wantString, tc.level.String())
	}
}

func TestLevelFromString(t *testing.T) {
	t.Parallel()

	type testCase struct {
		levelName string
		wantLevel log.Level
		wantOK    bool
	}

	testCases := []testCase{
		{"error", log.LevelError, true},
		{"warn", log.LevelWarn, true},
		{"info", log.LevelInfo, true},
		{"debug", log.LevelDebug, true},
		{"trace", log.LevelTrace, true},
		{"disabled", log.LevelDisabled, true},
		{"something-else", log.LevelUnknown, false},
	}

	for _, tc := range testCases {
		level, ok := log.LevelFromString(tc.levelName)

		assert.Equal(t, tc.wantLevel, level)
		assert.Equal(t, tc.wantOK, ok)
	}
}
