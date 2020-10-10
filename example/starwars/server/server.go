package main

import (
	"log"
	"os"

	"github.com/99designs/gqlgen/example/starwars"
	"github.com/99designs/gqlgen/example/starwars/generated"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/valyala/fasthttp"
)

const defaultPort = ":8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	gqlHandler := handler.NewDefaultServer(generated.NewExecutableSchema(starwars.NewResolver())).Handler()
	playground := playground.Handler("GraphQL playground", "/query")

	h := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/query":
			gqlHandler(ctx)
		case "/":
			playground(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(fasthttp.ListenAndServe(defaultPort, h))
}
