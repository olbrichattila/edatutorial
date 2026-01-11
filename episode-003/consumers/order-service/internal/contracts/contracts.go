package contracts

import "github.com/olbrichattila/edatutorial/shared/actions"

type OrderRepository interface {
	Save(ord actions.OrderCreatedAction) (int64, error)
}
