package v1resp

import (
	"time"

	"article-service/model"

	"github.com/google/uuid"
)

type CreateArticleDTO struct {
	ID uuid.UUID `json:"id"`
}

type ListArticlesDTO struct {
	RecordsCount int64        `json:"recordsCount"`
	Articles     []ArticleDTO `json:"articles"`
}

type ArticleDTO struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	Author    AuthorDTO `json:"author"`
}

func (dto *ArticleDTO) Convert(article *model.Article) ArticleDTO {
	respDto := ArticleDTO{
		ID:        article.ID,
		Title:     article.Title,
		Body:      article.Body,
		CreatedAt: article.CreatedAt,
		Author: AuthorDTO{
			ID:   article.Author.ID,
			Name: article.Author.Name,
		},
	}

	return respDto
}

func (dto *ListArticlesDTO) Convert(articles []*model.Article, recordsCount int64) ListArticlesDTO {
	responseDTO := ListArticlesDTO{
		RecordsCount: recordsCount,
	}

	for _, article := range articles {
		articleDTO := new(ArticleDTO).Convert(article)
		responseDTO.Articles = append(responseDTO.Articles, articleDTO)
	}

	return responseDTO
}
