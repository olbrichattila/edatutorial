package contracts

type LoggerRepository interface {
	Save(logType, actionID, correlationID, causationID, msg string, ind int64) error
}
