package initialization

import "go.uber.org/zap"

func InitializeClient() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Unable to initialize logger")
	}

	return logger, nil
}
