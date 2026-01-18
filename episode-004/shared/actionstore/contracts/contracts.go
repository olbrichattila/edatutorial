package contracts

import "github.com/olbrichattila/edatutorial/shared/actions"

type ActionStore interface {
	IsDuplicateAction(consumer string, actionMetaData actions.MetaData) (bool, string, error)
	SetDuplicateAction(consumer string, actionMetaData actions.MetaData, metadata string) error
	GetMetadata(consumer string, actionMetaData actions.MetaData) (string, error)
}
