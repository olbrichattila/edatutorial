package contracts

type LoggerRepository interface {
	Save(msg string) error
}
