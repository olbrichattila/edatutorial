package eventlogger

import (
	"github.com/olbrichattila/edatutorial/shared/actions"
	eventContracts "github.com/olbrichattila/edatutorial/shared/event/contracts"
	loggerContracts "github.com/olbrichattila/edatutorial/shared/logger/contracts"
)

const (
	logTopic = "logmessagecreated"
)

func New(evt eventContracts.EventManager) loggerContracts.Logger {
	return &logger{
		evt: evt,
	}
}

type logger struct {
	evt eventContracts.EventManager
}

func (l *logger) Info(msg string) error {
	return l.publish(actions.LogTypeInfo, msg)
}

func (l *logger) Error(msg string) error {
	return l.publish(actions.LogTypeError, msg)
}

func (l *logger) publish(logType actions.LogType, msg string) error {
	envelope := actions.New(actions.LogAction{LogType: logType, Message: msg})
	envelopeAsJSON, err := envelope.ToJSON()
	if err != nil {
		return err
	}
	return l.evt.Publish(logTopic, envelopeAsJSON)
}
