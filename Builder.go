package log

import (
	"strings"
)

type LogConfiguration interface {
	GetString(id string, defaultValue string) string
}

func CreateLogger(config LogConfiguration) Logger {
	minLevel := GetLoggingLevel(config.GetString("sw_logger_min_level", "trace"))

	switch config.GetString("sw_logger", "json") {
	case "json":
		return NewJsonLogger(minLevel)
	default:
		return NewConsoleLogger(minLevel)
	}

}

func GetLoggingLevel(minLevel string) LogLevel {
	minLevelStr := strings.TrimSpace(minLevel)

	if minLevelStr == "" {
		return LogLevelTrace
	}

	switch strings.ToLower(minLevelStr[0:1]) {
	case "d":
		return LogLevelDebug
	case "i":
		return LogLevelInfo
	case "w":
		return LogLevelWarn
	case "e":
		return LogLevelError
	default:
		return LogLevelTrace
	}
}
