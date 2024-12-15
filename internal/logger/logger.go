package logger

import "go.uber.org/zap"

var Log *zap.Logger = zap.NewNop()

func InitLogger() error {
	var err error
	Log, err = zap.NewProduction()
	if err != nil {
		return err
	}
	defer Log.Sync()
	return nil
}
