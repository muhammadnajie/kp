package router

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/muhammadnajie/kp/graph"
	"github.com/muhammadnajie/kp/graph/generated"
	"github.com/muhammadnajie/kp/internal/auth"
	"github.com/muhammadnajie/kp/internal/links"
	"github.com/muhammadnajie/kp/internal/users"
)

func Routes() chi.Router {
	router := chi.NewRouter()
	router.Use(auth.Middleware())

	router.Post("/login", users.Login)
	router.Post("/register", users.Register)

	router.Get("/links", links.GetAllController)
	router.Get("/links-by-title", links.GetByTitleController)
	router.Post("/link/create", links.CreateLinkController)
	router.Put("/link/update", links.UpdateLinkController)
	router.Delete("/link", links.DeleteLinkController)

	//===================== GRAPHQL
	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	router.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	router.Handle("/query", server)
	//===================== END OF GRAPHQL

	return router
}
