package log

// logging.go provides a simpler way to call LogInfoFields/LogErrorFields,
// while supporting ctx.

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"article-service/lib"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	_ "time/tzdata"
)

var logger = InitLog()

func InitLog() *logrus.Logger {
	log := logrus.New()
	filename := getAbsoluteLogPath()

	logFile := &lumberjack.Logger{
		Filename: filename, // Set your log file path
		MaxSize:  10,       // Maximum size of the log file before rotation (MB)
		Compress: false,    // Compress rotated log files
	}

	// Create a multi-writer to output logs to both terminal (stdout) and log file
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Set the log output to the multi-writer
	log.SetOutput(multiWriter)

	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.DateOnly,
		FullTimestamp:   true,
	})

	return log
}

// Builder builds up a log entry. Prefer using Infof or Errorf directly.
func Builder(ctx context.Context) *LogBuilder {
	return newBuilder(ctx)
}

func Infof(ctx context.Context, msg string, args ...any) {
	newBuilder(ctx).
		WithSource(lib.WhoCalledMe()).
		Now().Infof(msg, args...)
}

func Errorf(ctx context.Context, err error, msg string, args ...any) {
	newBuilder(ctx).
		WithSource(lib.WhoCalledMe()).
		WithError(err).
		Now().Errorf(msg, args...)
}

func getAbsoluteLogPath() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Cannot get executable path: %v", err)
	}

	// Assume logs should go relative to the project root
	projectRoot := filepath.Dir(exePath)

	logPath := filepath.Join(projectRoot, "storage", "log.txt")
	return logPath
}
