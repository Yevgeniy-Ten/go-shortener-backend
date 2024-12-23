package logger

import "go.uber.org/zap"

type MyLogger = *zap.Logger

func InitLogger() (MyLogger, error) {
	myLogger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	//nolint:errcheck // ignore error because it's not important
	defer myLogger.Sync()
	return myLogger, nil
}
