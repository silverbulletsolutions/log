package log

import (
	"context"
)

const logContextValueKey = "sw_gcm-l-ctx"

func FromContext(ctx context.Context) Logger {
	l, ok := ctx.Value(logContextValueKey).(Logger)
	if !ok {
		return nil
	}
	return l
}

func WithContext(parent context.Context, l Logger) context.Context {
	return context.WithValue(parent, logContextValueKey, l)
}

