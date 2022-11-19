package telegram_notifier

import (
	"os"
	"testing"
)

func TestNotify(t *testing.T) {
	_ = os.Setenv("TELEGRAM_TOKEN", "5808559800:AAEb2JxWD8V-69EHac_jXYUXmxOneLQOSKA")
	_ = os.Setenv("TELEGRAM_CHAT_ID", "876700915")

	Notify("google",
		"go-github",
		"https://github.com/google/go-github/releases/tag/v48.1.0",
		"v48.0.0", "v48.1.0", "minor", false)
}
