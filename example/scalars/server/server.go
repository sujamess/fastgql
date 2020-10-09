package main

import (
	"log"

	"github.com/99designs/gqlgen/example/scalars"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	playground := playground.Handler("Starwars", "/query")
	gqlHandler := handler.NewDefaultServer(scalars.NewExecutableSchema(scalars.Config{Resolvers: &scalars.Resolver{}})).Handler()

	app.All("/query", func(c *fiber.Ctx) error {
		gqlHandler(c.Context())
		return nil
	})

	app.All("/", func(c *fiber.Ctx) error {
		playground(c.Context())
		return nil
	})

	log.Fatal(app.Listen(":8084"))
}
