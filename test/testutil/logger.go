package testutil

import "share-basket-server/core/logger"

type dummyLogger struct{}

func (d *dummyLogger) Debug(msg string) {}
func (d *dummyLogger) Error(msg string) {}
func (d *dummyLogger) Info(msg string)  {}
func (d *dummyLogger) Warn(msg string)  {}
func (d *dummyLogger) With(key string, val any) logger.Logger {
	return d
}

func NewDummyLogger() logger.Logger {
	return &dummyLogger{}
}
