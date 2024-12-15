package logger

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

type CustomFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string
	// Colors defines whether to use colors in output
	Colors bool
}

func (f *CustomFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// Timestamp
	timestamp := entry.Time.Format(f.TimestampFormat)
	if f.Colors {
		fmt.Fprintf(b, "%s ", color.HiBlackString(timestamp))
	} else {
		fmt.Fprintf(b, "%s ", timestamp)
	}

	// Level
	var levelColor func(format string, a ...interface{}) string
	switch entry.Level {
	case log.DebugLevel:
		levelColor = color.HiBlackString
	case log.InfoLevel:
		levelColor = color.HiBlueString
	case log.WarnLevel:
		levelColor = color.YellowString
	case log.ErrorLevel, log.FatalLevel, log.PanicLevel:
		levelColor = color.RedString
	default:
		levelColor = color.WhiteString
	}

	if f.Colors {
		fmt.Fprintf(b, "[%s] ", levelColor(strings.ToUpper(entry.Level.String())))
	} else {
		fmt.Fprintf(b, "[%s] ", strings.ToUpper(entry.Level.String()))
	}

	// Function name and file
	if entry.HasCaller() {
		funcName := filepath.Base(entry.Caller.Function)
		fileName := filepath.Base(entry.Caller.File)
		if f.Colors {
			fmt.Fprintf(b, "%s:%d %s(): ",
				color.HiBlackString(fileName),
				entry.Caller.Line,
				color.CyanString(funcName),
			)
		} else {
			fmt.Fprintf(b, "%s:%d %s(): ", fileName, entry.Caller.Line, funcName)
		}
	}

	// Message
	if f.Colors {
		fmt.Fprintf(b, "%s ", color.WhiteString(entry.Message))
	} else {
		fmt.Fprintf(b, "%s ", entry.Message)
	}

	// Sort fields by key for consistent output
	var keys []string
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Fields
	if len(entry.Data) > 0 {
		if f.Colors {
			fmt.Fprint(b, color.HiBlackString("{"))
		} else {
			fmt.Fprint(b, "{")
		}
		for i, key := range keys {
			value := entry.Data[key]
			if f.Colors {
				fmt.Fprintf(b, "%s=%v", color.HiCyanString(key), value)
			} else {
				fmt.Fprintf(b, "%s=%v", key, value)
			}
			if i < len(keys)-1 {
				fmt.Fprint(b, " ")
			}
		}
		if f.Colors {
			fmt.Fprint(b, color.HiBlackString("}"))
		} else {
			fmt.Fprint(b, "}")
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

// InitLogger initializes the logger with custom formatter
func InitLogger(debug bool) {
	// Set logger level
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// Enable caller information
	log.SetReportCaller(true)

	// Set custom formatter
	log.SetFormatter(&CustomFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		Colors:          true,
	})
}
