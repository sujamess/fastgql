//go:generate go run ../../../testdata/gqlgen.go
package main

import (
	"log"
	"os"

	"github.com/sujamess/fastgql/example/federation/products/graph"
	"github.com/sujamess/fastgql/example/federation/products/graph/generated"
	"github.com/sujamess/fastgql/graphql/handler"
	"github.com/sujamess/fastgql/graphql/handler/debug"
	"github.com/sujamess/fastgql/graphql/playground"
	"github.com/valyala/fasthttp"
)

const defaultPort = "4002"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	srv.Use(&debug.Tracer{})

	playground := playground.Handler("GraphQL playground", "/query")
	gqlHandler := srv.Handler()

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
	log.Fatal(fasthttp.ListenAndServe(":"+port, h))
}
