package log

import (
	"go.uber.org/zap/zapcore"
)

//Client ... contains a production logger and development logger.
type Client struct {
	Logger ProdInterface
}

//ProdInterface ... contains all used zap log methods on prod.
//go:generate counterfeiter . ProdInterface
type ProdInterface interface {
	Info(msg string, fields ...zapcore.Field)
	Debug(msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
}

//New ... Creates a new instance of DB
func New(logInterface ProdInterface) *Client {
	return &Client{
		Logger: logInterface,
	}
}

//Info ... is the logger for the main info level (shown in production).
func (l *Client) Info(msg string, fields ...zapcore.Field) {
	l.Logger.Info(msg, fields...)
}

//Debug ... is the logger for the main debug level (hidden in production).
func (l *Client) Debug(msg string, fields ...zapcore.Field) {
	l.Logger.Debug(msg, fields...)
}

//Warn ... is the logger for the main warn level (shown in production).
func (l *Client) Warn(msg string, fields ...zapcore.Field) {
	l.Logger.Warn(msg, fields...)
}

//Error ... is the logger for the main error level (shown in production).
func (l *Client) Error(msg string, fields ...zapcore.Field) {
	l.Logger.Error(msg, fields...)
}
