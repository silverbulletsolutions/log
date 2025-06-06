package log

import (
	"encoding/json"
	"fmt"
	"strconv"
	"context"
)

type LogLevel int

const (
	LogLevelTrace LogLevel = 1
	LogLevelDebug LogLevel = 2
	LogLevelInfo  LogLevel = 3
	LogLevelWarn  LogLevel = 4
	LogLevelError LogLevel = 5
)

type LogValue interface {
	GetString() string
	GetJson() string
}

type Data map[string]LogValue

type Logger interface {
	Warn(message string, data Data)
	Info(message string, data Data)
	Debug(message string, data Data)
	Error(message string, data Data)
	Trace(message string, data Data)
	Panic(message string, panic error, data Data)

	Scope(scope Data) Logger
	ScopeData() Data

	AugmentWithContext(ctx context.Context) Logger
}

type Int int

func (l Int) GetString() string {
	return strconv.Itoa(int(l))
}

func (l Int) GetJson() string {
	return strconv.Itoa(int(l))
}

type UInt uint

func (l UInt) GetString() string {
	return strconv.FormatUint(uint64(l), 10)
}

func (l UInt) GetJson() string {
	return strconv.FormatUint(uint64(l), 10)
}

type Str string

func (l Str) GetString() string {
	return string(l)
}

func (l Str) GetJson() string {
	str := string(l)
	return jsonEscape(str)
}

func jsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func Err(err error) LogValue {
	return Str(err.Error())
}

type Bool bool

func (l Bool) GetString() string {
	if l {
		return "true"
	}
	return "false"
}

func (l Bool) GetJson() string {
	if l {
		return "true"
	}
	return "false"
}

type Float64 float64

func (l Float64) GetString() string {
	return fmt.Sprintf("%f", l)
}

func (l Float64) GetJson() string {
	return fmt.Sprintf("%f", l)
}
