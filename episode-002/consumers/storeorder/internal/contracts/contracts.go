package contracts

import "producer.example/internal/dto"

type OrderRepository interface {
	Save(ord dto.Order) (int64, error)
}
