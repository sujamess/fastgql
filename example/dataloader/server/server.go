package main

import (
	"log"

	"github.com/99designs/gqlgen/example/dataloader"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	playground := playground.Handler("Dataloader", "/query")
	gqlHandler := handler.NewDefaultServer(dataloader.NewExecutableSchema(dataloader.Config{Resolvers: &dataloader.Resolver{}})).Handler()

	app.Use(dataloader.LoaderMiddleware)

	app.All("/query", func(c *fiber.Ctx) error {
		gqlHandler(c.Context())
		return nil
	})

	app.All("/", func(c *fiber.Ctx) error {
		playground(c.Context())
		return nil
	})

	log.Println("connect to http://localhost:8082/ for graphql playground")
	log.Fatal(app.Listen(":8082"))
}
