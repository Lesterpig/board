package alert

import (
	"github.com/Lesterpig/board/probe"
	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

// Alerter is the interface all alerters should implement.
type Alerter interface {
	Alert(status probe.Status, category, serviceName, message, link, date string)
}
