package notify

import (
	"runtime"
	"strings"
)

const (
	SeverityLow = iota
	SeverityNormal
	SeverityUrgent
)

type Severity int

func (s Severity) String() string {
	sev := "low"

	switch s {
	case SeverityLow:
		sev = "low"
	case SeverityNormal:
		sev = "normal"
	case SeverityUrgent:
		sev = "critical"
	}

	if //goland:noinspection GoBoolExpressions
	runtime.GOOS == "darwin" {
		sev = strings.Title(sev)
	}
	if //goland:noinspection GoBoolExpressions
	runtime.GOOS == "windows" {
		switch s {
		case SeverityLow:
			sev = "Info"
		case SeverityNormal:
			sev = "Warning"
		case SeverityUrgent:
			sev = "Error"
		}
	}
	return sev
}

type Notify struct {
	title    string
	message  string
	severity Severity
}

func New(title, message string, severity Severity) *Notify {
	return &Notify{
		title:    title,
		message:  message,
		severity: severity,
	}
}
