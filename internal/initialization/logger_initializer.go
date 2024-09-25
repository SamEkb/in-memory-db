package initialization

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"in-memory-db/internal/configuration"
)

const defaultLvl = zapcore.InfoLevel
const defaultOutput = "./logs/output.log"

func InitializeLogger(conf *configuration.LoggingConfig) (*zap.Logger, error) {
	lvl := defaultLvl
	output := defaultOutput

	if conf != nil {
		if conf.Level != "" {
			ok := true
			lvl, ok = getLogLevel(conf.Level)
			if !ok {
				return nil, fmt.Errorf("log level: %s wasn't found", conf.Level)
			}

			if conf.Output != "" {
				output = conf.Output
				logDir := filepath.Dir(conf.Output)
				if _, err := os.Stat(logDir); os.IsNotExist(err) {
					err = os.MkdirAll(logDir, os.ModePerm)
					if err != nil {
						return nil, fmt.Errorf("failed to create log directory: %v", err)
					}
				}
			}
		}
	}

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(lvl),
		OutputPaths: []string{output},
		Encoding:    "json",
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, err
}

func getLogLevel(level string) (zapcore.Level, bool) {
	levels := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
	}

	lvl, ok := levels[level]
	return lvl, ok
}
