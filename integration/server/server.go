package main

import (
	"context"
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/valyala/fasthttp"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/integration"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	cfg := integration.Config{Resolvers: &integration.Resolver{}}
	cfg.Complexity.Query.Complexity = func(childComplexity, value int) int {
		// Allow the integration client to dictate the complexity, to verify this
		// function is executed.
		return value
	}

	srv := handler.NewDefaultServer(integration.NewExecutableSchema(cfg))
	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		if e, ok := errors.Cause(e).(*integration.CustomError); ok {
			return &gqlerror.Error{
				Message: e.UserMessage,
				Path:    graphql.GetFieldContext(ctx).Path(),
			}
		}
		return graphql.DefaultErrorPresenter(ctx, e)
	})
	srv.Use(extension.FixedComplexityLimit(1000))

	serverHandler := srv.Handler()
	playgroundHandler := playground.Handler("GraphQL playground", "/query")

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/":
			playgroundHandler(ctx)
		case "/query":
			serverHandler(ctx)
		default:
			ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		}
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(fasthttp.ListenAndServe(":"+port, requestHandler))
}
