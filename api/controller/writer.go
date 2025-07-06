package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"article-service/api/apiconst"
	"article-service/dto/response"
	"article-service/infrastructure/log"
)

func WriteSuccess(ctx context.Context, w http.ResponseWriter, statusCode int, result interface{}) {
	err := writeJSON(w, statusCode, response.SuccessResponse{
		Success: true,
		Result:  result,
	})
	if err != nil {
		log.Errorf(ctx, err, "error while writing JSON response")
	}
}

func WriteError(ctx context.Context, w http.ResponseWriter, statusCode int, err error) {
	err = writeJSON(w, statusCode, response.FailureResponse{
		Success: false,
		Failure: err.Error(),
	})
	if err != nil {
		log.Errorf(ctx, err, "error while writing JSON response")
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, response interface{}) error {
	w.Header().Set(apiconst.ContentTypeHeader, apiconst.ContentTypeJSON)
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(response)
}
