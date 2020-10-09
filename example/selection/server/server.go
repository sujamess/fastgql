package main

import (
	"log"

	"github.com/99designs/gqlgen/example/selection"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	playground := playground.Handler("Selection Demo", "/query")
	gqlHandler := handler.NewDefaultServer(selection.NewExecutableSchema(selection.Config{Resolvers: &selection.Resolver{}})).Handler()

	app.All("/query", func(c *fiber.Ctx) error {
		gqlHandler(c.Context())
		return nil
	})

	app.All("/", func(c *fiber.Ctx) error {
		playground(c.Context())
		return nil
	})

	log.Fatal(app.Listen(":8086"))
}
