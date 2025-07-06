package elasticsearch

import (
	"article-service/configloader"
	"article-service/infrastructure/log"
	"context"

	"github.com/olivere/elastic/v7"
)

type ElasticSearch struct {
	Client *elastic.Client
}

var elasticInstance *ElasticSearch

func GetElasticInstance() ElasticSearch {
	return *elasticInstance
}

func InitElasticSearch(ctx context.Context, config configloader.ElasticConfig) {
	es, err := elastic.NewClient(elastic.SetURL(config.URL), elastic.SetSniff(false))
	if err != nil {
		log.Errorf(ctx, err, "[ElasticSearch] Failed to connect")
		panic(err.Error())
	}

	elasticInstance = &ElasticSearch{
		Client: es,
	}
	log.Infof(ctx, "[ElasticSearch] Initializing instance")
}

func InitElasticSearchMock() {
	elasticInstance = &ElasticSearch{
		Client: nil,
	}
}
