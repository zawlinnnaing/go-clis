package notify

import (
	"fmt"
	"os/exec"
)

var command = exec.Command

func (n Notify) Send() error {
	cmdName := "terminal-notifier"
	cmd, err := exec.LookPath(cmdName)
	if err != nil {
		return err
	}
	title := fmt.Sprintf("(%s) %s", n.severity.String(), n.title)
	notifyCommand := command(cmd, "-title", title, "-message", n.message)
	return notifyCommand.Run()
}
