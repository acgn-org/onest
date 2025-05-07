package logfield

import (
	log "github.com/sirupsen/logrus"
)

func New(com string) LogWithFields {
	logWithFields := LogWithFields{
		components: com,
	}
	logWithFields.FieldLogger = logWithFields.NewLogger()
	return logWithFields
}

type LogWithFields struct {
	log.FieldLogger

	components string
	action     string
}

func (l LogWithFields) NewLogger() log.FieldLogger {
	logger := log.WithField(Component, l.components)
	if l.action != "" {
		logger = logger.WithField(Action, l.action)
	}
	return logger
}

func (l LogWithFields) WithComponent(component string) LogWithFields {
	l.components += ":" + component
	l.FieldLogger = l.NewLogger()
	return l
}

func (l LogWithFields) WithAction(action string) LogWithFields {
	l.action = action
	l.FieldLogger = l.NewLogger()
	return l
}
