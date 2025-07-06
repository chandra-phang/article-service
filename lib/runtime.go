package lib

import (
	"fmt"
	"runtime"
	"strings"
)

func WhoCalledMe() string {
	_, fileName, lineNo, ok := runtime.Caller(2)
	if !ok {
		return "(failed to get caller)"
	}
	fileName = trimpath(fileName)

	callerSourceFile := fmt.Sprintf("%s:%d", fileName, lineNo)
	return callerSourceFile
}

func trimpath(path string) string {
	substringStart := "article-service"

	// Find the index of the substring
	index := strings.Index(path, substringStart)

	// Check if the substring is found in the string
	if index != -1 {
		// Extract the substring starting from the found index to the end
		result := path[index:]
		path = result
	}

	return path
}
