package eventlogger

import (
	"github.com/olbrichattila/edatutorial/shared/actions"
	eventContracts "github.com/olbrichattila/edatutorial/shared/event/contracts"
	loggerContracts "github.com/olbrichattila/edatutorial/shared/logger/contracts"
)

const (
	logEventTopic = "log-event-emitted"
	consumer      = "eventlogger"
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
	return l.publish(actions.LogTypeInfo, msg, "", "", "", 0)
}

func (l *logger) InfoWithAction(e actions.MetaData, msg string) error {
	return l.publish(
		actions.LogTypeInfo,
		msg,
		e.ActionID,
		e.CorrelationID,
		e.CausationID,
		e.Index,
	)
}

func (l *logger) ErrorWithAction(e actions.MetaData, msg string) error {
	return l.publish(
		actions.LogTypeError,
		msg,
		e.ActionID,
		e.CorrelationID,
		e.CausationID,
		e.Index,
	)
}

func (l *logger) Error(msg string) error {
	return l.publish(actions.LogTypeError, msg, "", "", "", 0)
}

func (l *logger) publish(logType actions.LogType, msg, actionID, correlationID, causationID string, index int64) error {
	envelope := actions.New(logEventTopic, actions.LogAction{
		LogType:       logType,
		ActionID:      actionID,
		CorrelationID: correlationID,
		CausationID:   causationID,
		Index:         index,
		Message:       msg,
	})
	envelopeAsJSON, err := envelope.ToJSON()
	if err != nil {
		return err
	}
	return l.evt.Publish(logEventTopic, envelopeAsJSON)
}
