// Package log contains custom logger for our API
package log // import "github.com/la4ezar/restapi/pkg/log

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type logKey struct{}

var (
	formatters = map[string]logrus.Formatter{
		"text": &logrus.TextFormatter{},
		"json": &logrus.JSONFormatter{},
	}

	outputs = map[string]io.Writer{
		os.Stdout.Name(): os.Stdout,
		os.Stderr.Name(): os.Stderr,
	}

	mutex = sync.RWMutex{}

	D = DefaultLogger
	C = LoggerFromContext
)

func Configure(ctx context.Context, cfg *Config) (context.Context, error) {
	mutex.Lock()
	defer mutex.Unlock()

	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	formatter := formatters[cfg.Format]
	output := outputs[cfg.Output]

	entry := logrus.NewEntry(logrus.StandardLogger())
	entry.Logger.SetLevel(level)
	entry.Logger.SetFormatter(formatter)
	entry.Logger.SetOutput(output)

	return ContextWithLogger(ctx, entry), nil
}

func ContextWithLogger(ctx context.Context, entry *logrus.Entry) context.Context {
	return context.WithValue(ctx, logKey{}, entry)
}

func LoggerFromContext(ctx context.Context) *logrus.Entry {
	mutex.RLock()
	defer mutex.RUnlock()

	entry := ctx.Value(logKey{})
	if entry == nil {
		entry = logrus.NewEntry(logrus.StandardLogger())
	}

	return copyFromEntry(entry.(*logrus.Entry))
}

func DefaultLogger() *logrus.Entry {
	return LoggerFromContext(context.Background())
}

func RegisterFormatter(formatterName string, formatter logrus.Formatter) error {
	if _, exists := formatters[formatterName]; exists {
		return fmt.Errorf("formatter with name %s is already registered", formatterName)
	}
	formatters[formatterName] = formatter
	return nil
}

func RequestLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ContextWithLogger(r.Context(), DefaultLogger())
			r = r.WithContext(ctx)

			start := time.Now()

			remoteAddr := r.RemoteAddr
			if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
				remoteAddr = realIP
			}

			beforeLogger := D().WithFields(logrus.Fields{
				"request": r.RequestURI,
				"method":  r.Method,
				"remote":  remoteAddr,
			})

			beforeLogger.Info("Starting handling request...")

			next.ServeHTTP(w, r)

			duration := time.Since(start)

			afterLogger := D().WithFields(logrus.Fields{
				"status_code": http.StatusOK,
				"took":        duration,
			})

			afterLogger.Info("Finished handling request...")
		})
	}
}

func copyFromEntry(entry *logrus.Entry) *logrus.Entry {
	entryData := make(logrus.Fields, len(entry.Data))
	for k, v := range entry.Data {
		entryData[k] = v
	}
	newEntry := logrus.NewEntry(entry.Logger)
	newEntry.Data = entryData
	newEntry.Time = entry.Time
	newEntry.Level = entry.Level
	newEntry.Buffer = entry.Buffer
	newEntry.Message = entry.Message

	return newEntry
}
