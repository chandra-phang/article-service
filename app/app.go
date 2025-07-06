package app

import (
	"context"
	"time"

	"article-service/api"
	"article-service/application"
	"article-service/configloader"
	"article-service/db/db_client"
	"article-service/infrastructure/elasticsearch"
	"article-service/infrastructure/log"
)

type Application struct {
}

// Returns a new instance of the application
func NewApplication() Application {
	return Application{}
}

func (a Application) InitApplication(configFilePath string) {
	time.Local = time.UTC
	ctx := context.Background()
	log.Infof(ctx, "[App] Application is starting up")

	if err := configloader.LoadConfigFromFile(configFilePath); err != nil {
		log.Errorf(ctx, err, "[App] failed to load config, path: %s", configFilePath)
		panic(err)
	}

	config := configloader.GetRootConfig()

	a.initDB(ctx, config.DbConfig)
	a.initElasticSearch(ctx, config.ElasticConfig)
	a.initServices()
	a.initRoutes(config.AppConfig)
}

func (a Application) initDB(ctx context.Context, cfg configloader.DbConfig) {
	db_client.InitDatabase(ctx, cfg)
	db_client.RunMigrations(ctx, cfg)
}

func (a Application) initElasticSearch(ctx context.Context, cfg configloader.ElasticConfig) {
	elasticsearch.InitElasticSearch(ctx, cfg)
}

func (a Application) initServices() {
	application.InitServices()
}

func (a Application) initRoutes(cfg configloader.AppConfig) {
	api.InitRoutes(cfg)
}
