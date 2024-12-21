package logger

import "go.uber.org/zap"

type MyLogger = *zap.Logger

func InitLogger() (MyLogger, error) {
	myLogger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	//nolint:errcheck // игнорируем ошибку потому что так принято
	defer myLogger.Sync()
	return myLogger, nil
}
