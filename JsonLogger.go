package log

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"context"
)

type JsonLogger struct {
	r        *regexp.Regexp
	minLevel LogLevel
	scope    Data
}

func NewJsonLogger(minLogLevel LogLevel) *JsonLogger {
	cl := JsonLogger{
		r:        regexp.MustCompile(`{[a-zA-Z_][a-zA-Z0-9_]*}`),
		minLevel: minLogLevel,
		scope:    make(Data),
	}

	return &cl
}

func (c *JsonLogger) Scope(newData Data) Logger {
	newScope := make(Data, len(c.scope)+len(newData))

	for k, v := range c.scope {
		newScope[k] = v
	}

	for k, v := range newData {
		newScope[k] = v
	}

	return &JsonLogger{
		r:        c.r,
		minLevel: c.minLevel,
		scope:    newScope,
	}
}

func (c *JsonLogger) AugmentWithContext(ctx context.Context) Logger {
	source := FromContext(ctx)
	if source == nil {
		return c
	}

	newData := source.ScopeData()
	if len(newData) == 0 {
		return c
	}

	newScope := make(Data, len(c.scope)+len(newData))

	for k, v := range newData {
		newScope[k] = v
	}

	for k, v := range c.scope {
		newScope[k] = v
	}

	return &JsonLogger{
		r:        c.r,
		minLevel: c.minLevel,
		scope:    newScope,
	}
}

func (c *JsonLogger) ScopeData() Data {
	return c.scope
}

func (c *JsonLogger) Panic(message string, panicError error, data Data) {
	_, ok := data["error"]
	if !ok {
		data["error"] = Err(panicError)
	}
	c.writeLog(LogLevelError, message, data)
	panic(panicError)
}

func (c *JsonLogger) Warn(message string, data Data) {
	c.writeLog(LogLevelWarn, message, data)
}

func (c *JsonLogger) Info(message string, data Data) {
	c.writeLog(LogLevelInfo, message, data)
}

func (c *JsonLogger) Debug(message string, data Data) {
	c.writeLog(LogLevelDebug, message, data)
}

func (c *JsonLogger) Error(message string, data Data) {
	c.writeLog(LogLevelError, message, data)
}

func (c *JsonLogger) Trace(message string, data Data) {
	c.writeLog(LogLevelTrace, message, data)
}

func (c *JsonLogger) writeLog(ll LogLevel, message string, data Data) {
	if ll < c.minLevel {
		return
	}

	var out strings.Builder

	filledMessage := c.fillMessage(message, data)

	out.WriteString("{")
	out.WriteString("\"timestamp\":\"")
	out.WriteString(timestamp())
	out.WriteString("\",\"level\":\"")
	out.WriteString(getLevelString(ll))
	out.WriteString("\",\"message\":")
	out.WriteString(Str(filledMessage).GetJson())
	for k, v := range c.scope {
		if _, inData := data[k]; inData {
			continue
		}
		out.WriteString(",\"")
		out.WriteString(k)
		out.WriteString("\":")
		out.WriteString(v.GetJson())
	}
	for k, v := range data {
		if strings.HasPrefix(k, "_") {
			continue
		}
		out.WriteString(",\"")
		out.WriteString(k)
		out.WriteString("\":")
		out.WriteString(v.GetJson())
	}

	out.WriteString("}")

	fmt.Println(out.String())
}

func timestamp() string {
	t := time.Now().UTC()
	return t.Format("2006-01-02T15:04:05.000Z")
}

func (c *JsonLogger) fillMessage(message string, data Data) string {
	return c.r.ReplaceAllStringFunc(message, func(match string) string {
		field := match[1 : len(match)-1]
		d, ok := data[field]
		if ok {
			return d.GetString()
		}

		if field == "" || field[0] == '_' {
			return match
		}

		// scope data can be used in string templates
		d, ok = c.scope[field]
		if ok {
			return d.GetString()
		}

		return match
	})
}

func getLevelString(level LogLevel) string {
	switch level {
	case LogLevelTrace:
		return "Trace"
	case LogLevelDebug:
		return "Debug"
	case LogLevelInfo:
		return "Info"
	case LogLevelWarn:
		return "Warn"
	case LogLevelError:
		return "Error"
	}
	return "?"
}
