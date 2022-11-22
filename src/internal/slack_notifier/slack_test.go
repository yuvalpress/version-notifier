package slack_notifier

import (
	"os"
	"testing"
)

func TestNotify(t *testing.T) {
	// TODO: Fill this values using GitHub Actions
	_ = os.Setenv("SLACK_TOKEN", "xoxb-4424054153393-4424080095233-hebAA0UQ2D0oX4GlhCkBpXds")
	_ = os.Setenv("SLACK_CHANNEL", "C04C0ED6JER")

	Notify("google",
		"go-github",
		"https://github.com/google/go-github/releases/tag/v48.1.0",
		"v48.0.0", "v48.1.0", "minor", false)
}
