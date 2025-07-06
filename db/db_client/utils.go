package db_client

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func normalizeWhitespace(s string) string {
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	return s
}

func generateTransactionID() string {
	n := time.Now().UnixNano()
	base36 := strconv.FormatInt(n, 36)
	base36 = strings.ToUpper(base36)
	// trim the leading 5 chars, since they're the most-significant bits that are mostly the same
	return fmt.Sprintf("txnID::%s", base36[5:])
}
