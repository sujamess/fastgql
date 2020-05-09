package main

import (
	"log"
	"os"

	"github.com/99designs/gqlgen/example/starwars"
	"github.com/99designs/gqlgen/example/starwars/generated"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	app := fiber.New()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(starwars.NewResolver()))
	serverHandler := srv.Handler()
	playgroundHandler := playground.Handler("GraphQL playground", "/query")

	app.Use("/", func(c *fiber.Ctx) {
		playgroundHandler(c.Fasthttp)
	})

	app.Use("/query", func(c *fiber.Ctx) {
		serverHandler(c.Fasthttp)
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	app.Listen(defaultPort)
}
