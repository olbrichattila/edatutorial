package contracts

type Logger interface {
	Info(msg string) error
	Error(msg string) error
}
