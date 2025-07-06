package model

import (
	"time"

	"github.com/google/uuid"
)

const ArticleIndex = "articles"

type Article struct {
	ID        uuid.UUID
	Title     string
	Body      string
	CreatedAt time.Time
	Author    Author
}
