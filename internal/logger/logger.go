package logger

import "go.uber.org/zap"

var Log *zap.Logger = zap.NewNop()

func InitLogger() error {
	var err error
	Log, err = zap.NewProduction()
	if err != nil {
		return err
	}
	//nolint:errcheck // игнорируем ошибку потому что так принято
	defer Log.Sync()
	return nil
}
