package contracts

import "github.com/olbrichattila/edatutorial/shared/actions"

type OrderRepository interface {
	Save(ord actions.OrderSentAction) (int64, error)
}
