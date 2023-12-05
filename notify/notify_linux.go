package notify

import "os/exec"

var command = exec.Command

func (n Notify) Send() error {
	notifyCmdName := "notify-send"
	notifyCmd, err := exec.LookPath(notifyCmdName)
	if err != nil {
		return err
	}
	notifyCommand := command(notifyCmd, "-u", n.severity.String(), n.title, n.message)
	return notifyCommand.Run()
}
