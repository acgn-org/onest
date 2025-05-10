package logfield

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&nested.Formatter{
		TimestampFormat: "01/02 15:04:05",
		FieldsOrder: []string{
			Component, Action,
		},
	})
}

func New(com string) LoggerWithFields {
	logWithFields := LoggerWithFields{
		Entry:      log.WithField(Component, com),
		components: com,
	}
	return logWithFields
}

type LoggerWithFields struct {
	*log.Entry

	components string
	action     string
}

func (l LoggerWithFields) WithAction(action string) LoggerWithFields {
	l.action = action
	l.Entry = log.WithFields(log.Fields{
		Component: l.components,
		Action:    l.action,
	})
	return l
}
