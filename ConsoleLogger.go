package log

import (
	"fmt"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"
	"context"
)

type ConsoleColors struct {
	Reset  string
	Red    string
	Green  string
	Yellow string
	Blue   string
	Purple string
	Cyan   string
	Gray   string
	White  string
}

func GetConsoleColors() ConsoleColors {
	if runtime.GOOS == "windows" {
		return ConsoleColors{}
	}
	return ConsoleColors{
		Reset:  "\033[0m",
		Red:    "\033[31m",
		Green:  "\033[32m",
		Yellow: "\033[33m",
		Blue:   "\033[34m",
		Purple: "\033[35m",
		Cyan:   "\033[36m",
		Gray:   "\033[37m",
		White:  "\033[97m",
	}
}

type ConsoleLogger struct {
	r        *regexp.Regexp
	minLevel LogLevel
	colors   ConsoleColors
}

func NewConsoleLogger(minLogLevel LogLevel) *ConsoleLogger {

	cl := ConsoleLogger{
		r:        regexp.MustCompile(`{[a-zA-Z_][a-zA-Z0-9_]*}`),
		minLevel: minLogLevel,
		colors:   GetConsoleColors(),
	}

	return &cl
}

func (c *ConsoleLogger) Scope(scope Data) Logger {
	return c
}

func (c *ConsoleLogger) AugmentWithContext(ctx context.Context) Logger {
	return c
}

func (c *ConsoleLogger) ScopeData() Data {
	return Data{}
}

func (c *ConsoleLogger) Panic(message string, panicError error, data Data) {
	_, ok := data["error"]
	if !ok {
		data["error"] = Err(panicError)
	}
	c.writeLog(LogLevelError, message, data)
	panic(panicError)
}

func (c *ConsoleLogger) Warn(message string, data Data) {
	c.writeLog(LogLevelWarn, message, data)
}

func (c *ConsoleLogger) Info(message string, data Data) {
	c.writeLog(LogLevelInfo, message, data)
}

func (c *ConsoleLogger) Debug(message string, data Data) {
	c.writeLog(LogLevelDebug, message, data)
}

func (c *ConsoleLogger) Error(message string, data Data) {
	c.writeLog(LogLevelError, message, data)
}

func (c *ConsoleLogger) Trace(message string, data Data) {
	c.writeLog(LogLevelTrace, message, data)
}

func (c *ConsoleLogger) writeLog(ll LogLevel, message string, data Data) {
	if ll < c.minLevel {
		return
	}

	color := c.colors.Reset
	switch ll {
	case LogLevelError:
		color = c.colors.Red
	case LogLevelWarn:
		color = c.colors.Yellow
	case LogLevelInfo:
		color = c.colors.White
	case LogLevelTrace:
		color = c.colors.Gray
	}

	fmt.Printf(color+"%s %s %s%s\n"+c.colors.Reset, consoleTimestamp(), consoleLevel(ll), c.fillMessage(message, data), printData(data))
}

func consoleTimestamp() string {
	t := time.Now().UTC()
	return t.Format("2006-01-02 15:04:05.000")
}

func (c *ConsoleLogger) fillMessage(message string, data Data) string {
	return c.r.ReplaceAllStringFunc(message, func(match string) string {
		field := match[1 : len(match)-1]
		d, ok := data[field]
		if !ok {
			return match
		}
		return d.GetString()
	})
}

func printData(data Data) string {
	sb := strings.Builder{}

	keys := make([]string, 0, len(data))
	for k := range data {
		if k[0] == '_' {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := data[k]
		sb.WriteString(fmt.Sprintf(" %s=%s", k, v.GetString()))
	}
	return sb.String()
}

func consoleLevel(level LogLevel) string {
	switch level {
	case LogLevelTrace:
		return "T"
	case LogLevelDebug:
		return "D"
	case LogLevelInfo:
		return "I"
	case LogLevelWarn:
		return "W"
	case LogLevelError:
		return "E"
	}
	return "?"
}
