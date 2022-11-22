package telegram_notifier

import (
	"testing"
)

func TestSlackNotify(t *testing.T) {
	Notify("google",
		"go-github",
		"https://github.com/google/go-github/releases/tag/v48.1.0",
		"v48.0.0", "v48.1.0", "minor", false)
}

func TestSlackNotifyWithRelease(t *testing.T) {
	Notify("google",
		"go-github",
		"https://github.com/google/go-github/releases/tag/v48.1.0",
		"v48.0.0", "v48.1.0", "minor", true)
}
