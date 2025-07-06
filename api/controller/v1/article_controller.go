package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"article-service/api/controller"
	"article-service/apperror"
	"article-service/application"
	v1req "article-service/dto/request/v1_req"
	v1resp "article-service/dto/response/v1_resp"
	"article-service/infrastructure/log"
)

type articleController struct {
	svc application.IArticleService
}

func InitArticleController() *articleController {
	return &articleController{
		svc: application.GetArticleService(),
	}
}

func (c articleController) CreateArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqBody, _ := io.ReadAll(r.Body)
	dto := v1req.CreateArticleDTO{}
	if err := json.Unmarshal(reqBody, &dto); err != nil {
		log.Errorf(ctx, err, "[V1][ArticleController][CreateArticle] Failed to unmarshal request body %v into dto", reqBody)
		controller.WriteError(ctx, w, http.StatusBadRequest, apperror.ErrUnmarshalRequestBodyFailed)
		return
	}

	err := dto.Validate(ctx)
	if err != nil {
		log.Errorf(ctx, err, "[V1][ArticleController][CreateArticle] Validation failed for request dto %v ", dto)
		controller.WriteError(ctx, w, http.StatusBadRequest, err)
		return
	}

	id, err := c.svc.CreateArticle(ctx, dto)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == apperror.ErrAuthorNotFound {
			statusCode = http.StatusUnprocessableEntity
		}
		log.Errorf(ctx, err, "[V1][ArticleController][CreateArticle] svc.CreateArticle is failed for request dto: %v ", dto)
		controller.WriteError(ctx, w, statusCode, err)
		return
	}

	resp := v1resp.CreateArticleDTO{ID: id}
	controller.WriteSuccess(ctx, w, http.StatusCreated, resp)
}

func (c articleController) ListArticles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queryParams := r.URL.Query()

	query := queryParams.Get("query")
	authorName := queryParams.Get("authorName")
	sortBy := queryParams.Get("sortBy")
	sortDirection := queryParams.Get("sortDirection")
	limit := queryParams.Get("limit")
	page := queryParams.Get("page")

	limitInt, _ := strconv.Atoi(limit)
	pageInt, _ := strconv.Atoi(page)

	dto := v1req.ListArticlesDTO{
		Query:         query,
		AuthorName:    authorName,
		SortBy:        sortBy,
		SortDirection: sortDirection,
		Limit:         limitInt,
		Page:          pageInt,
	}

	err := dto.Validate(ctx)
	if err != nil {
		log.Errorf(ctx, err, "[V1][ArticleController][ListArticles] Validation failed for request dto %v ", dto)
		controller.WriteError(ctx, w, http.StatusBadRequest, err)
		return
	}

	articles, recordsCount, err := c.svc.ListArticles(ctx, dto)
	if err != nil {
		log.Errorf(ctx, err, "[V1][ArticleController][ListArticles] svc.ListArticles is failed")
		controller.WriteError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	resp := new(v1resp.ListArticlesDTO).Convert(articles, recordsCount)
	controller.WriteSuccess(ctx, w, http.StatusOK, resp)
}
