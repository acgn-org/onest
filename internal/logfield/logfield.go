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
		components: com,
	}
	logWithFields.Entry = logWithFields.NewLogger()
	return logWithFields
}

type LoggerWithFields struct {
	*log.Entry

	components string
	action     string
}

func (l LoggerWithFields) NewLogger() *log.Entry {
	logger := log.WithField(Component, l.components)
	if l.action != "" {
		logger = logger.WithField(Action, l.action)
	}
	return logger
}

func (l LoggerWithFields) WithSubComponent(component string) LoggerWithFields {
	l.components += ":" + component
	l.Entry = l.NewLogger()
	return l
}

func (l LoggerWithFields) WithAction(action string) LoggerWithFields {
	l.action = action
	l.Entry = l.NewLogger()
	return l
}
