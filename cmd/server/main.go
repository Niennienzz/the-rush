package main

import (
	"fmt"
	"log"
	"net/http"

	graphqlHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"the-rush/graph"
	"the-rush/graph/runtime"
	httpHandler "the-rush/handler"
	"the-rush/middleware"
	"the-rush/repository"
)

func main() {
	var (
		port     = "8080"
		repo     = repository.NewLocal()
		gHandler = graphqlHandler.NewDefaultServer(runtime.NewExecutableSchema(runtime.Config{Resolvers: graph.NewResolver(repo)}))
		dHandler = httpHandler.NewDownloadHandler(repo)
	)

	http.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))
	http.Handle("/graphql", middleware.CORS(gHandler))
	http.Handle("/download/csv", middleware.CORS(dHandler))

	log.Printf("connect to http://localhost:%s/playground to try the GraphQL API\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
