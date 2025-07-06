package factory

import (
	"article-service/model"
	"log"
	"time"

	"github.com/google/uuid"
)

var SampleAuthorChandra model.Author
var SampleAuthorPhang model.Author
var SampleArticle1 model.Article
var SampleArticle2 model.Article

func init() {
	SampleAuthorChandra = model.Author{
		ID:   uuid.MustParse("0197da8f-47ed-78b1-7b0f-ea4f4a1af25e"),
		Name: "Chandra",
	}
	SampleAuthorPhang = model.Author{
		ID:   uuid.MustParse("0197da8f-47ed-78b1-7b0f-ea4f4a1af25f"),
		Name: "Phang",
	}

	parsedTime1, err := time.Parse(time.RFC3339, "2025-07-05T09:00:00+07:00")
	if err != nil {
		log.Fatal(err)
	}

	SampleArticle1 = model.Article{
		ID:        uuid.MustParse("0197db1c-c6c4-7140-bee3-8efd703f30c8"),
		Title:     "Satu satu aku sayang ibu",
		Body:      "Dua dua juga sayang ayah",
		CreatedAt: parsedTime1,
		Author:    SampleAuthorChandra,
	}

	parsedTime2, err := time.Parse(time.RFC3339, "2025-07-05T10:00:00+07:00")
	if err != nil {
		log.Fatal(err)
	}

	SampleArticle2 = model.Article{
		ID:        uuid.MustParse("0197db1c-c6c4-7140-bee3-8efd703f30c9"),
		Title:     "Tiga tiga sayang adik kakak",
		Body:      "Satu dua tiga, sayang semuanya",
		CreatedAt: parsedTime2,
		Author:    SampleAuthorPhang,
	}
}
