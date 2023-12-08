//go:build !integration

package notify

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		s Severity
	}{
		{SeverityLow}, {SeverityNormal}, {SeverityUrgent},
	}
	for _, testCase := range testCases {
		expTitle := "Title"
		expMessage := "Message"
		notification := New(expTitle, expMessage, testCase.s)
		if notification.severity != testCase.s {
			t.Errorf("Expected severity: %s, recevied: %s", testCase.s, notification.severity)
		}
		if notification.title != expTitle {
			t.Errorf("Expected title: %s, recevied: %s", expTitle, notification.title)
		}
		if notification.message != expMessage {
			t.Errorf("Expected message: %s, recevied: %s", expMessage, notification.message)
		}
	}
}

func TestSeverityString(t *testing.T) {
	testCases := []struct {
		severity Severity
		os       string
		exp      string
	}{
		{
			severity: SeverityLow,
			os:       "linux",
			exp:      "low",
		},
		{
			severity: SeverityNormal,
			os:       "linux",
			exp:      "normal",
		},
		{
			severity: SeverityUrgent,
			os:       "linux",
			exp:      "critical",
		},
		{
			severity: SeverityLow,
			os:       "windows",
			exp:      "Info",
		},
		{
			severity: SeverityNormal,
			os:       "windows",
			exp:      "Warning",
		},
		{
			severity: SeverityUrgent,
			os:       "windows",
			exp:      "Error",
		},
		{
			severity: SeverityLow,
			os:       "darwin",
			exp:      "Low",
		},
		{
			severity: SeverityNormal,
			os:       "darwin",
			exp:      "Normal",
		},
		{
			severity: SeverityUrgent,
			os:       "darwin",
			exp:      "Critical",
		},
	}
	for _, testCase := range testCases {
		name := fmt.Sprintf("%s(%s)", testCase.os, testCase.severity)
		t.Run(name, func(t *testing.T) {
			if testCase.os != runtime.GOOS {
				t.Skip("Skipped: not OS:", testCase.os)
			}
			sev := testCase.severity.String()
			if sev != testCase.exp {
				t.Errorf("Expected severity string: %s, recevied: %s", testCase.exp, sev)
			}
		})
	}
}

func mockCmd(exe string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess"}
	cs = append(cs, exe)
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER=1"}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER") != "1" {
		return
	}
	cmdName := ""
	switch runtime.GOOS {
	case "linux":
		cmdName = "notify-send"
	case "darwin":
		cmdName = "terminal-notifier"
	case "windows":
		cmdName = "powershell"
	}
	if strings.Contains(os.Args[2], cmdName) {
		os.Exit(0)
	}
	os.Exit(1)
}

func TestSend(t *testing.T) {
	notification := New("test title", "test message", SeverityLow)
	command = mockCmd
	err := notification.Send()
	if err != nil {
		t.Error(err)
	}
}
