package middleware

import (
	"fmt"
	"net/http"
	"time"

	"article-service/infrastructure/log"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

func LogRequest(next http.Handler) http.Handler {
	return middleware.RequestLogger(logMW{})(next)
}

type logMW struct{}
type logWriter struct {
	req   *http.Request
	entry *log.LogBuilder
}

func (m logMW) NewLogEntry(r *http.Request) middleware.LogEntry {
	ctx := r.Context()
	fields := getFieldsFromRequest(r)
	e := log.Builder(ctx).WithFields(fields)
	return logWriter{r, e}
}

func (w logWriter) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	w.entry.WithFields(logrus.Fields{
		"resp_status": status,
	}).Now().Infof("HTTP request completed")
}

func (w logWriter) Panic(p interface{}, stack []byte) {
	w.entry.WithFields(logrus.Fields{
		"stack": string(stack),
	}).Now().Errorf("panic: %s", fmt.Sprint(p))
}

func getFieldsFromRequest(r *http.Request) map[string]any {
	logFields := map[string]any{
		"http_method": r.Method,
		"uri":         r.RequestURI,
	}
	return logFields
}
