//go:build !containers
// +build !containers

package app

import "github.com/zawlinnnaing/go-clis/notify"

func sendNotification(msg string) {
	notification := notify.New("Pomodoro CLI", msg, notify.SeverityNormal)

	notification.Send()
}
