package main

import (
	"log"
	"os"

	"github.com/99designs/gqlgen/example/starwars"
	"github.com/99designs/gqlgen/example/starwars/generated"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
)

const defaultPort = ":8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	app := fiber.New()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(starwars.NewResolver()))
	serverHandler := srv.Handler()
	playgroundHandler := playground.Handler("GraphQL playground", "/query")

	app.All("/query", func(c *fiber.Ctx) error {
		serverHandler(c.Context())
		return nil
	})

	app.All("/", func(c *fiber.Ctx) error {
		playgroundHandler(c.Context())
		return nil
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	log.Fatal(app.Listen(defaultPort))
}
