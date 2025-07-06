// Package appctx provides operations for interacting with ctx.Context objects
// that are specific to this application, such as defining our custom context key.
// This package is named 'appctx' to differentiate from 'ctx' which is the
// conventional naming of a context.Context object.
package appctx

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ctxKey is a custom type for context keys.
// All context keys should be private to this package.
// Reading and writing to the context should be done through this package.
type ctxKey string

// Context keys
const (
	requestIDKey = ctxKey("request ID")
)

func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		return reqID
	}
	return ""
}

func WithReqID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, requestIDKey, reqID)
}

func GenerateRequestID() string {
	n := time.Now().UnixNano()
	base36 := strconv.FormatInt(n, 36)
	base36 = strings.ToUpper(base36)
	// trim the leading 5 chars, since they're the most-significant bits that are mostly the same
	return fmt.Sprintf("reqID::%s", base36)
}
