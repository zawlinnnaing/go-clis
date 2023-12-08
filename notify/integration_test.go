//go:build integration
// +build integration

package notify_test

import (
	"github.com/zawlinnnaing/go-clis/notify"
	"testing"
)

func TestSend(t *testing.T) {
	notification := notify.New("Test title", "Test message", notify.SeverityLow)
	err := notification.Send()
	if err != nil {
		t.Error(err)
	}
}
