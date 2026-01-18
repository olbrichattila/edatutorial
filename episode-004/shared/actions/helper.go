package actions

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type MetaData struct {
	ActionType    string    `json:"action_type"`
	ActionID      string    `json:"action_id"`
	CausationID   string    `json:"causation_id"`
	CorrelationID string    `json:"correlation_id"`
	Index         int64     `json:"index"`
	OccurredAt    time.Time `json:"occurred_at"`
}

// Episode 004 Added causation and correlation ID
type Envelope[T any] struct {
	MetaData MetaData `json:"metadata"`
	Payload  T        `json:"payload"`
}

func New[T any](actionType string, payload T) Envelope[T] {
	actionID := uuid.NewString()
	return Envelope[T]{
		MetaData: MetaData{
			ActionType:    actionType,
			ActionID:      actionID,
			CorrelationID: actionID,
			OccurredAt:    time.Now().UTC(),
		},
		Payload: payload,
	}
}

func NewFromParent[TP, T any](actionType string, parentEnvelope Envelope[TP], index int64, payload T) Envelope[T] {
	return Envelope[T]{
		MetaData: MetaData{
			ActionType:    actionType,
			ActionID:      uuid.NewString(),
			CorrelationID: parentEnvelope.MetaData.CorrelationID,
			CausationID:   parentEnvelope.MetaData.ActionID,
			Index:         index,
			OccurredAt:    time.Now().UTC(),
		},

		Payload: payload,
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
