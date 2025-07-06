package search

import (
	"article-service/apperror"
	"article-service/infrastructure/elasticsearch"
	"article-service/infrastructure/log"
	"article-service/model"
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/olivere/elastic/v7"
)

type ArticleSearch struct {
	Client *elastic.Client
}

type ArticleSearchDoc struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Body  string    `json:"body"`
}

func GetArticleSearch() IArticleSearch {
	instance := elasticsearch.GetElasticInstance()
	return ArticleSearch{
		Client: instance.Client,
	}
}

func (s ArticleSearch) Index(ctx context.Context, article model.Article) error {
	doc := ArticleSearchDoc{
		ID:    article.ID,
		Title: article.Title,
		Body:  article.Body,
	}
	_, err := s.Client.Index().
		Index(model.ArticleIndex).
		Id(article.ID.String()).
		BodyJson(doc).
		Do(ctx)

	if err != nil {
		log.Errorf(ctx, err, "[ArticleSearch][Index] Index is failed, index: %s, doc: %v", model.ArticleIndex, doc)
		return apperror.ErrIndexElasticFailed
	}

	return nil
}

func (s ArticleSearch) Search(ctx context.Context, query string) ([]uuid.UUID, error) {
	res, err := s.Client.Search().
		Index(model.ArticleIndex).
		Query(elastic.NewMultiMatchQuery(query, "title", "body")).
		Sort("_score", false).
		Do(ctx)

	if err != nil {
		log.Errorf(ctx, err, "[ArticleSearch][Search] Search is failed, query: %s", query)
		return nil, apperror.ErrSearchElasticFailed
	}

	var ids []uuid.UUID
	for _, hit := range res.Hits.Hits {
		var a model.Article
		json.Unmarshal(hit.Source, &a)
		ids = append(ids, a.ID)
	}

	return ids, nil
}
