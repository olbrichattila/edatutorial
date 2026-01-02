package contracts

type LoggerRepository interface {
	Save(logType, actionId, msg string) error
}
