package contracts

type LoggerRepository interface {
	Save(logType, actionID, msg string) error
}
