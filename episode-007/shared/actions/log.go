package actions

const (
	LogTypeInfo  LogType = "info"
	LogTypeError LogType = "error"
)

type LogType string

type LogAction struct {
	ActionID      string  `json:"actionID"`
	CausationID   string  `json:"causationID"`
	CorrelationID string  `json:"correlationID"`
	Index         int64   `json:"Index"`
	LogType       LogType `json:"logType"`
	Message       string  `json:"message"`
}
