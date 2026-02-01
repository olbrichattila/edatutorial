package contracts

import "github.com/olbrichattila/edatutorial/shared/actions"

type Logger interface {
	Info(msg string) error
	InfoWithAction(e actions.MetaData, msg string) error
	Error(msg string) error
	ErrorWithAction(e actions.MetaData, msg string) error
}
