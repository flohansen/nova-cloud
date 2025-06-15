package logging

type Logger interface {
	Error(msg string, v ...any)
	Info(msg string, v ...any)
}
