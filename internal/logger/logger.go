package logger

import "go.uber.org/zap"

type MyLogger interface {
	Error(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
}

func InitLogger() (MyLogger, error) {
	myLogger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	//nolint:errcheck // ignore error because it's not important
	defer myLogger.Sync()
	return myLogger, nil
}
