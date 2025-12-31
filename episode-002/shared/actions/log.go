package actions

const (
	LogTypeInfo  LogType = "info"
	LogTypeError LogType = "error"
)

type LogType string

type LogAction struct {
	LogType LogType `json:"logType"`
	Message string  `json:"string"`
}
