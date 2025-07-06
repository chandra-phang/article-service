package tests

import (
	"article-service/api/apiconst"
	v1 "article-service/api/controller/v1"
	"article-service/application"
	"article-service/configloader"
	"article-service/db/db_client"
	v1req "article-service/dto/request/v1_req"
	"article-service/dto/response"
	v1resp "article-service/dto/response/v1_resp"
	"article-service/factory"
	"article-service/infrastructure/elasticsearch"
	"article-service/infrastructure/log"
	"article-service/search"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var config configloader.RootConfig

func TestMain(m *testing.M) {
	ctx := context.Background()

	if err := configloader.LoadConfigFromFile("../config.yml"); err != nil {
		log.Errorf(ctx, err, "[IntegrationTest][TestMain] failed to load config")
		panic(fmt.Sprintf("[IntegrationTest][TestMain] failed to load config, err=%v", err))
	}

	config = configloader.GetRootConfig()

	elasticsearch.InitElasticSearch(context.Background(), config.ElasticConfig)

	search.GetArticleSearch().Index(context.Background(), factory.SampleArticle1)
	search.GetArticleSearch().Index(context.Background(), factory.SampleArticle2)

	db_client.DropTestDB(ctx, config.DbConfig)
	db_client.CreateTestDB(ctx, config.DbConfig)
	defer db_client.DropTestDB(ctx, config.DbConfig)

	m.Run()
}

func TestCreateArticle(t *testing.T) {
	db_client.RunIntegrationTestSeed(context.Background(), config.DbConfig, "../tests/seed")

	dto := v1req.CreateArticleDTO{
		Title:    "New article title",
		Body:     "New article body",
		AuthorId: "0197da8f-47ed-78b1-7b0f-ea4f4a1af25e",
	}
	reqBody, _ := json.Marshal(dto)
	r := httptest.NewRequest(http.MethodPost, "/v1/articles", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	application.InitArticleService()
	articleController := v1.InitArticleController()
	articleController.CreateArticle(w, r)

	res := w.Result()
	defer res.Body.Close()

	respBody := response.SuccessResponse{}
	respBytes, _ := io.ReadAll(w.Body)
	json.Unmarshal(respBytes, &respBody)

	result := v1resp.CreateArticleDTO{}
	resultJSON, _ := json.Marshal(result)
	json.Unmarshal(resultJSON, &result)

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, true, respBody.Success)
	assert.NotNil(t, result.ID)

	db_client.TruncateTestDB(context.Background(), config.DbConfig)
}

func TestListArticles_WithNoQuery_ReturnAllRecords(t *testing.T) {
	db_client.RunIntegrationTestSeed(context.Background(), config.DbConfig, "../tests/seed")

	r := httptest.NewRequest(http.MethodPost, "/v1/articles", nil)
	w := httptest.NewRecorder()

	r.URL = &url.URL{}
	r.URL.RawQuery = ""
	r.Header.Add(apiconst.ContentTypeHeader, apiconst.ContentTypeJSON)

	application.InitArticleService()
	articleController := v1.InitArticleController()
	articleController.ListArticles(w, r)

	res := w.Result()
	defer res.Body.Close()

	respBody := response.SuccessResponse{}
	respBytes, _ := io.ReadAll(w.Body)
	json.Unmarshal(respBytes, &respBody)

	result := v1resp.ListArticlesDTO{}
	resultJSON, _ := json.Marshal(respBody.Result)
	json.Unmarshal(resultJSON, &result)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, true, respBody.Success)
	assert.Equal(t, 2, len(result.Articles))
	assert.Equal(t, int64(2), result.RecordsCount)

	article := result.Articles[0]
	expArticle2 := factory.SampleArticle2

	assert.Equal(t, expArticle2.ID, article.ID)
	assert.Equal(t, expArticle2.Title, article.Title)
	assert.Equal(t, expArticle2.Body, article.Body)
	assert.Equal(t, expArticle2.CreatedAt, article.CreatedAt)
	assert.Equal(t, expArticle2.Author.ID, article.Author.ID)
	assert.Equal(t, expArticle2.Author.Name, article.Author.Name)

	article2 := result.Articles[1]
	expArticle := factory.SampleArticle1

	assert.Equal(t, expArticle.ID, article2.ID)
	assert.Equal(t, expArticle.Title, article2.Title)
	assert.Equal(t, expArticle.Body, article2.Body)
	assert.Equal(t, expArticle.CreatedAt, article2.CreatedAt)
	assert.Equal(t, expArticle.Author.ID, article2.Author.ID)
	assert.Equal(t, expArticle.Author.Name, article2.Author.Name)

	db_client.TruncateTestDB(context.Background(), config.DbConfig)
}

func TestListArticles_WithQuery_ReturnAllRecords(t *testing.T) {
	db_client.RunIntegrationTestSeed(context.Background(), config.DbConfig, "../tests/seed")

	r := httptest.NewRequest(http.MethodPost, "/v1/articles", nil)
	w := httptest.NewRecorder()

	r.URL = &url.URL{}
	query := r.URL.Query()
	query.Add("query", "Satu")
	r.URL.RawQuery = query.Encode()
	r.Header.Add(apiconst.ContentTypeHeader, apiconst.ContentTypeJSON)

	application.InitArticleService()
	articleController := v1.InitArticleController()
	articleController.ListArticles(w, r)

	res := w.Result()
	defer res.Body.Close()

	respBody := response.SuccessResponse{}
	respBytes, _ := io.ReadAll(w.Body)
	json.Unmarshal(respBytes, &respBody)

	result := v1resp.ListArticlesDTO{}
	resultJSON, _ := json.Marshal(respBody.Result)
	json.Unmarshal(resultJSON, &result)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, true, respBody.Success)
	assert.Equal(t, 2, len(result.Articles))
	assert.Equal(t, int64(2), result.RecordsCount)

	article := result.Articles[0]
	expArticle2 := factory.SampleArticle2

	assert.Equal(t, expArticle2.ID, article.ID)
	assert.Equal(t, expArticle2.Title, article.Title)
	assert.Equal(t, expArticle2.Body, article.Body)
	assert.Equal(t, expArticle2.CreatedAt, article.CreatedAt)
	assert.Equal(t, expArticle2.Author.ID, article.Author.ID)
	assert.Equal(t, expArticle2.Author.Name, article.Author.Name)

	article2 := result.Articles[1]
	expArticle := factory.SampleArticle1

	assert.Equal(t, expArticle.ID, article2.ID)
	assert.Equal(t, expArticle.Title, article2.Title)
	assert.Equal(t, expArticle.Body, article2.Body)
	assert.Equal(t, expArticle.CreatedAt, article2.CreatedAt)
	assert.Equal(t, expArticle.Author.ID, article2.Author.ID)
	assert.Equal(t, expArticle.Author.Name, article2.Author.Name)

	db_client.TruncateTestDB(context.Background(), config.DbConfig)
}

func TestListArticles_WithQuery_ReturnOneRecord(t *testing.T) {
	db_client.RunIntegrationTestSeed(context.Background(), config.DbConfig, "../tests/seed")

	r := httptest.NewRequest(http.MethodPost, "/v1/articles", nil)
	w := httptest.NewRecorder()

	r.URL = &url.URL{}
	query := r.URL.Query()
	query.Add("query", "ayah")
	r.URL.RawQuery = query.Encode()
	r.Header.Add(apiconst.ContentTypeHeader, apiconst.ContentTypeJSON)

	application.InitArticleService()
	articleController := v1.InitArticleController()
	articleController.ListArticles(w, r)

	res := w.Result()
	defer res.Body.Close()

	respBody := response.SuccessResponse{}
	respBytes, _ := io.ReadAll(w.Body)
	json.Unmarshal(respBytes, &respBody)

	result := v1resp.ListArticlesDTO{}
	resultJSON, _ := json.Marshal(respBody.Result)
	json.Unmarshal(resultJSON, &result)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, true, respBody.Success)
	assert.Equal(t, 1, len(result.Articles))
	assert.Equal(t, int64(1), result.RecordsCount)

	article2 := result.Articles[0]
	expArticle := factory.SampleArticle1

	assert.Equal(t, expArticle.ID, article2.ID)
	assert.Equal(t, expArticle.Title, article2.Title)
	assert.Equal(t, expArticle.Body, article2.Body)
	assert.Equal(t, expArticle.CreatedAt, article2.CreatedAt)
	assert.Equal(t, expArticle.Author.ID, article2.Author.ID)
	assert.Equal(t, expArticle.Author.Name, article2.Author.Name)

	db_client.TruncateTestDB(context.Background(), config.DbConfig)
}

func TestListArticles_WithAuthorName_ReturnOneRecord(t *testing.T) {
	db_client.RunIntegrationTestSeed(context.Background(), config.DbConfig, "../tests/seed")

	r := httptest.NewRequest(http.MethodPost, "/v1/articles", nil)
	w := httptest.NewRecorder()

	r.URL = &url.URL{}
	query := r.URL.Query()
	query.Add("authorName", "Phang")
	r.URL.RawQuery = query.Encode()
	r.Header.Add(apiconst.ContentTypeHeader, apiconst.ContentTypeJSON)

	application.InitArticleService()
	articleController := v1.InitArticleController()
	articleController.ListArticles(w, r)

	res := w.Result()
	defer res.Body.Close()

	respBody := response.SuccessResponse{}
	respBytes, _ := io.ReadAll(w.Body)
	json.Unmarshal(respBytes, &respBody)

	result := v1resp.ListArticlesDTO{}
	resultJSON, _ := json.Marshal(respBody.Result)
	json.Unmarshal(resultJSON, &result)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, true, respBody.Success)
	assert.Equal(t, 1, len(result.Articles))
	assert.Equal(t, int64(1), result.RecordsCount)

	article2 := result.Articles[0]
	expArticle := factory.SampleArticle2

	assert.Equal(t, expArticle.ID, article2.ID)
	assert.Equal(t, expArticle.Title, article2.Title)
	assert.Equal(t, expArticle.Body, article2.Body)
	assert.Equal(t, expArticle.CreatedAt, article2.CreatedAt)
	assert.Equal(t, expArticle.Author.ID, article2.Author.ID)
	assert.Equal(t, expArticle.Author.Name, article2.Author.Name)

	db_client.TruncateTestDB(context.Background(), config.DbConfig)
}
