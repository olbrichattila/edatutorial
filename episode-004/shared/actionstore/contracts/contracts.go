package contracts

import "github.com/olbrichattila/edatutorial/shared/actions"

type ActionStore interface {
	IsDuplicateAction(actionMetaData actions.MetaData) (bool, error)
}
