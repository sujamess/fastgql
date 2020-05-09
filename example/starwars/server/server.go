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

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(starwars.NewResolver()))
	// srv.AroundFields(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	// 	rc := graphql.GetFieldContext(ctx)
	// 	fmt.Println("Entered", rc.Object, rc.Field.Name)
	// 	res, err = next(ctx)
	// 	fmt.Println("Left", rc.Object, rc.Field.Name, "=>", res, err)
	// 	return res, err
	// })

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
