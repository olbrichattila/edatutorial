package actions

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Envelope[T any] struct {
	ActionID   string    `json:"action_id"`
	OccurredAt time.Time `json:"occurred_at"`
	Payload    T         `json:"payload"`
}

func New[T any](payload T) Envelope[T] {
	id := uuid.NewString()

	return Envelope[T]{
		ActionID:   id,
		OccurredAt: time.Now().UTC(),
		Payload:    payload,
	}
}

func (e Envelope[T]) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func FromJSON[T any](data []byte) (Envelope[T], error) {
	var env Envelope[T]
	err := json.Unmarshal(data, &env)
	return env, err
}
