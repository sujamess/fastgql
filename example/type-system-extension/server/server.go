package main

import (
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/valyala/fasthttp"

	extension "github.com/99designs/gqlgen/example/type-system-extension"
	"github.com/99designs/gqlgen/graphql/handler"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	playground := playground.Handler("GraphQL playground", "/query")
	gqlHandler := handler.NewDefaultServer(
		extension.NewExecutableSchema(
			extension.Config{
				Resolvers: extension.NewRootResolver(),
				Directives: extension.DirectiveRoot{
					EnumLogging:   extension.EnumLogging,
					FieldLogging:  extension.FieldLogging,
					InputLogging:  extension.InputLogging,
					ObjectLogging: extension.ObjectLogging,
					ScalarLogging: extension.ScalarLogging,
					UnionLogging:  extension.UnionLogging,
				},
			},
		),
	).Handler()

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
