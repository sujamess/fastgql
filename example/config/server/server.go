package main

import (
	"log"

	todo "github.com/99designs/gqlgen/example/config"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	playground := playground.Handler("Todo", "/query")
	gqlHandler := handler.NewDefaultServer(todo.NewExecutableSchema(todo.New())).Handler()

	app.All("/query", func(c *fiber.Ctx) error {
		gqlHandler(c.Context())
		return nil
	})

	app.All("/", func(c *fiber.Ctx) error {
		playground(c.Context())
		return nil
	})

	log.Fatal(app.Listen(":8081"))
}
