package logger

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/ony-boom/swoosh/config"
)

type logger struct {
	inner *log.Logger
}

var (
	instance *logger
	once     sync.Once
	Log      *logger = getLogger()
)

func getLogger() *logger {
	once.Do(func() {
		logFilePath := filepath.Join(config.BasePath(), "swoosh.log")
		f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		instance = &logger{inner: log.New(f, "", log.LstdFlags|log.Lshortfile)}
	})
	return instance
}

func (l *logger) Info(format string, args ...any) {
	l.inner.Printf("[INFO]: "+format, args...)
}

func (l *logger) Error(format string, args ...any) {
	l.inner.Printf("[ERROR]: "+format, args...)
}

func (l *logger) Debug(format string, args ...any) {
	l.inner.Printf("[DEBUG]: "+format, args...)
}

func (l *logger) Warn(format string, args ...any) {
	l.inner.Printf("[WARN]: "+format, args...)
}
