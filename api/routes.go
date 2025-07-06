package api

import (
	"net/http"

	v1 "article-service/api/controller/v1"
	"article-service/api/middleware"
	"article-service/configloader"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func newRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.SetRequestID)
	r.Use(middleware.LogRequest)

	return r
}

func InitRoutes(cfg configloader.AppConfig) {
	r := newRouter()

	r.Route("/v1", func(r chi.Router) {
		articleController := v1.InitArticleController()

		r.Route("/articles", func(r chi.Router) {
			r.Get("/", articleController.ListArticles)
			r.Post("/", articleController.CreateArticle)
		})
	})

	http.ListenAndServe(cfg.Port, r)
}
