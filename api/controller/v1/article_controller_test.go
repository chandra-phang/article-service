package v1

import (
	"article-service/api/apiconst"
	"article-service/apperror"
	"article-service/application"
	"article-service/application/mock_application"
	v1req "article-service/dto/request/v1_req"
	"article-service/dto/response"
	v1resp "article-service/dto/response/v1_resp"
	"article-service/factory"
	"article-service/infrastructure/elasticsearch"
	"article-service/model"
	"article-service/utils"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var createArticleDTO = v1req.CreateArticleDTO{
	Title:    "Article title",
	Body:     "Article body",
	AuthorId: "0197da8f-47ed-78b1-7b0f-ea4f4a1af25e",
}

func Test_InitArticleController(t *testing.T) {
	elasticsearch.InitElasticSearchMock()
	application.InitArticleService()

	articleController := InitArticleController()
	assert.NotNil(t, articleController.svc)
}

func Test_CreateArticle_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	dto := createArticleDTO

	articleID := utils.GenerateUUID()
	svc := mock_application.NewMockIArticleService(ctrl)
	svc.EXPECT().CreateArticle(gomock.Any(), dto).Return(articleID, nil)

	w := httptest.NewRecorder()
	r := &http.Request{Header: http.Header{}}

	jsonBytes, _ := json.Marshal(dto)
	r.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
	r.Header.Add(apiconst.ContentTypeHeader, apiconst.ContentTypeJSON)

	articleController{svc}.CreateArticle(w, r)
	statusCode := w.Result().StatusCode

	respBody := response.SuccessResponse{}
	respBytes, _ := io.ReadAll(w.Body)
	json.Unmarshal(respBytes, &respBody)

	result, _ := json.Marshal(respBody.Result)
	resultDTO := v1resp.CreateArticleDTO{}
	json.Unmarshal(result, &resultDTO)

	assert.Equal(t, http.StatusCreated, statusCode)
	assert.True(t, respBody.Success)
	assert.Equal(t, articleID, resultDTO.ID)
}

func Test_CreateArticle_ReturnErr_WhenInvalidJson(t *testing.T) {
	ctrl := gomock.NewController(t)

	svc := mock_application.NewMockIArticleService(ctrl)

	w := httptest.NewRecorder()
	r := &http.Request{Header: http.Header{}}

	jsonBody := "{"
	r.Body = io.NopCloser(bytes.NewBuffer([]byte(jsonBody)))
	r.Header.Add(apiconst.ContentTypeHeader, apiconst.ContentTypeJSON)

	articleController{svc}.CreateArticle(w, r)

	statusCode := w.Result().StatusCode
	respBytes, _ := io.ReadAll(w.Body)

	respBody := response.FailureResponse{}
	json.Unmarshal(respBytes, &respBody)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.False(t, respBody.Success)
	assert.Equal(t, apperror.ErrUnmarshalRequestBodyFailed.Error(), respBody.Failure)
}

func Test_CreateArticle_ReturnErr_WhenInvalidDTO(t *testing.T) {
	ctrl := gomock.NewController(t)
	dto := v1req.CreateArticleDTO{}

	svc := mock_application.NewMockIArticleService(ctrl)

	w := httptest.NewRecorder()
	r := &http.Request{Header: http.Header{}}

	jsonBytes, _ := json.Marshal(dto)
	r.Body = io.NopCloser(bytes.NewBuffer([]byte(jsonBytes)))
	r.Header.Add(apiconst.ContentTypeHeader, apiconst.ContentTypeJSON)

	articleController{svc}.CreateArticle(w, r)

	statusCode := w.Result().StatusCode
	respBytes, _ := io.ReadAll(w.Body)

	respBody := response.FailureResponse{}
	json.Unmarshal(respBytes, &respBody)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.False(t, respBody.Success)
	assert.Equal(t, "title is required, body is required, authorId is required", respBody.Failure)
}

func Test_CreateArticle_ReturnErr_WhenAuthorNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	dto := createArticleDTO

	svc := mock_application.NewMockIArticleService(ctrl)
	svc.EXPECT().CreateArticle(gomock.Any(), dto).Return(uuid.Nil, apperror.ErrAuthorNotFound)

	w := httptest.NewRecorder()
	r := &http.Request{Header: http.Header{}}

	jsonBytes, _ := json.Marshal(dto)
	r.Body = io.NopCloser(bytes.NewBuffer([]byte(jsonBytes)))
	r.Header.Add(apiconst.ContentTypeHeader, apiconst.ContentTypeJSON)

	articleController{svc}.CreateArticle(w, r)

	statusCode := w.Result().StatusCode
	respBytes, _ := io.ReadAll(w.Body)

	respBody := response.FailureResponse{}
	json.Unmarshal(respBytes, &respBody)

	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.False(t, respBody.Success)
	assert.Equal(t, apperror.ErrAuthorNotFound.Error(), respBody.Failure)
}

func Test_ListArticles_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	dto := v1req.ListArticlesDTO{}

	mockResult := []*model.Article{&factory.SampleArticle1}
	recordsCount := int64(1)
	svc := mock_application.NewMockIArticleService(ctrl)
	svc.EXPECT().ListArticles(gomock.Any(), dto).Return(mockResult, recordsCount, nil)

	w := httptest.NewRecorder()
	r := &http.Request{Header: http.Header{}}

	r.URL = &url.URL{}
	r.URL.RawQuery = ""
	r.Header.Add(apiconst.ContentTypeHeader, apiconst.ContentTypeJSON)

	articleController{svc}.ListArticles(w, r)
	statusCode := w.Result().StatusCode

	respBody := response.SuccessResponse{}
	respBytes, _ := io.ReadAll(w.Body)
	json.Unmarshal(respBytes, &respBody)

	result, _ := json.Marshal(respBody.Result)
	resultDTO := v1resp.ListArticlesDTO{}
	json.Unmarshal(result, &resultDTO)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.True(t, respBody.Success)
	assert.Equal(t, recordsCount, resultDTO.RecordsCount)

	assert.Equal(t, 1, len(resultDTO.Articles))
	assert.Equal(t, mockResult[0].ID, resultDTO.Articles[0].ID)
	assert.Equal(t, mockResult[0].Title, resultDTO.Articles[0].Title)
	assert.Equal(t, mockResult[0].Body, resultDTO.Articles[0].Body)
	assert.True(t, mockResult[0].CreatedAt.Equal(resultDTO.Articles[0].CreatedAt))
	assert.Equal(t, mockResult[0].Author.ID, resultDTO.Articles[0].Author.ID)
	assert.Equal(t, mockResult[0].Author.Name, resultDTO.Articles[0].Author.Name)
}

func Test_ListArticles_ReturnErr_WhenInvalidDTO(t *testing.T) {
	ctrl := gomock.NewController(t)

	svc := mock_application.NewMockIArticleService(ctrl)

	w := httptest.NewRecorder()
	r := &http.Request{Header: http.Header{}}

	r.URL = &url.URL{}
	query := r.URL.Query()
	query.Add("sortBy", "author_id")
	r.URL.RawQuery = query.Encode()
	r.Header.Add(apiconst.ContentTypeHeader, apiconst.ContentTypeJSON)

	articleController{svc}.ListArticles(w, r)
	statusCode := w.Result().StatusCode

	respBody := response.FailureResponse{}
	respBytes, _ := io.ReadAll(w.Body)
	json.Unmarshal(respBytes, &respBody)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.False(t, respBody.Success)
	assert.Equal(t, "sortBy should be one of created_at title author_name", respBody.Failure)
}

func Test_ListArticles_ReturnErr_WhenListArticlesFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	dto := v1req.ListArticlesDTO{}

	recordsCount := int64(0)
	svc := mock_application.NewMockIArticleService(ctrl)
	svc.EXPECT().ListArticles(gomock.Any(), dto).Return(nil, recordsCount, apperror.ErrGetRecordFailed)

	w := httptest.NewRecorder()
	r := &http.Request{Header: http.Header{}}

	r.URL = &url.URL{}
	r.URL.RawQuery = ""
	r.Header.Add(apiconst.ContentTypeHeader, apiconst.ContentTypeJSON)

	articleController{svc}.ListArticles(w, r)
	statusCode := w.Result().StatusCode

	respBody := response.FailureResponse{}
	respBytes, _ := io.ReadAll(w.Body)
	json.Unmarshal(respBytes, &respBody)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.False(t, respBody.Success)
	assert.Equal(t, apperror.ErrGetRecordFailed.Error(), respBody.Failure)
}
