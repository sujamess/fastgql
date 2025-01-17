# gqlgen with fasthttp

![gqlgen](https://user-images.githubusercontent.com/46195831/89802919-0bb8ef00-db2a-11ea-8ba4-88e7a58b2fd2.png)

## What is gqlgen?

[gqlgen](https://github.com/sujamess/fastgql) is a Go library for building GraphQL servers without any fuss.<br/>

- **gqlgen is based on a Schema first approach** — You get to Define your API using the GraphQL [Schema Definition Language](http://graphql.org/learn/schema/).
- **gqlgen prioritizes Type safety** — You should never see `map[string]interface{}` here.
- **gqlgen enables Codegen** — We generate the boring bits, so you can focus on building your app quickly.

Still not convinced enough to use **gqlgen**? Compare **gqlgen** with other Go graphql [implementations](https://gqlgen.com/feature-comparison/)

## Getting Started
- To install gqlgen run the command `go get github.com/sujamess/fastgql` in your project directory.<br/>
- Add this line to your go mod file `replace github.com/sujamess/fastgql v0.13.0 => github.com/sujamess/fastgql v0.13.2`
- You could initialize a new project using the recommended folder structure by running this command `go run github.com/sujamess/fastgql init`.

First work your way through the [Getting Started](https://gqlgen.com/getting-started/) tutorial.

If you can't find what your looking for, look at our [examples](https://github.com/sujamess/fastgql/tree/master/example) for example usage of gqlgen, or visit [godoc](https://godoc.org/github.com/sujamess/fastgql).

## Using with Fasthttp

```go
package main

import (
  "log"
  "github.com/valyala/fasthttp"
  "github.com/sujamess/fastgql"
  "<Your go module>/gql"
  "<Your go module>/generated"
  "github.com/sujamess/fastgql/graphql/handler"
  "github.com/sujamess/fastgql/graphql/playground"
)

func main() {
	playground := playground.Handler("GraphQL playground", "/query")
	gqlHandler := handler.NewDefaultServer(generated.NewExecutableSchema(gql.NewResolver())).Handler()

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

  log.Fatal(fasthttp.ListenAndServe(":8080", h))
}
```

## Using with Fiber

```go
package main

import (
  "log"
  "github.com/gofiber/fiber/v2"
  "github.com/sujamess/fastgql"
  "<Your go module>/gql"
  "<Your go module>/generated"
  "github.com/sujamess/fastgql/graphql/handler"
  "github.com/sujamess/fastgql/graphql/playground"
)

func main() {
  app := fiber.New()

  srv := handler.NewDefaultServer(generated.NewExecutableSchema(gql.NewResolver()))
	gqlHandler := srv.Handler()
  playground := playground.Handler("GraphQL playground", "/query")

  app.All("/query", func(c *fiber.Ctx) error {
    gqlHandler(c.Context())
    return nil
  })

  app.All("/playground", func(c *fiber.Ctx) error {
    playground(c.Context())
    return nil
  })

  log.Fatal(app.Listen(":8080"))
}
```
You could find a more comprehensive guide to help you get started [here](https://gqlgen.com/getting-started/).<br/>
We also have a couple of real-world [examples](https://github.com/sujamess/fastgql/tree/master/example) that show how to GraphQL applications with **gqlgen** seamlessly,
You can see these [examples](https://github.com/sujamess/fastgql/tree/master/example) here or visit [godoc](https://godoc.org/github.com/sujamess/fastgql).

## Reporting Issues

If you think you've found a bug, or something isn't behaving the way you think it should, please raise an [issue](https://github.com/sujamess/fastgql/issues) on GitHub.

## Contributing

We welcome contributions, Read our [Contribution Guidelines](https://github.com/sujamess/fastgql/blob/master/CONTRIBUTING.md) to learn more about contributing to **gqlgen**
## Frequently asked questions

### How do I prevent fetching child objects that might not be used?

When you have nested or recursive schema like this:

```graphql
type User {
  id: ID!
  name: String!
  friends: [User!]!
}
```

You need to tell gqlgen that it should only fetch friends if the user requested it. There are two ways to do this;

- #### Using Custom Models

Write a custom model that omits the friends field:

```go
type User struct {
  ID int
  Name string
}
```

And reference the model in `gqlgen.yml`:

```yaml
# gqlgen.yml
models:
  User:
    model: github.com/you/pkg/model.User # go import path to the User struct above
```

- #### Using Explicit Resolvers

If you want to Keep using the generated model, mark the field as requiring a resolver explicitly in `gqlgen.yml` like this:

```yaml
# gqlgen.yml
models:
  User:
    fields:
      friends:
        resolver: true # force a resolver to be generated
```

After doing either of the above and running generate we will need to provide a resolver for friends:

```go
func (r *userResolver) Friends(ctx context.Context, obj *User) ([]*User, error) {
  // select * from user where friendid = obj.ID
  return friends,  nil
}
```

### Can I change the type of the ID from type String to Type Int?

Yes! You can by remapping it in config as seen below:

```yaml
models:
  ID: # The GraphQL type ID is backed by
    model:
      - github.com/sujamess/fastgql/graphql.IntID # An go integer
      - github.com/sujamess/fastgql/graphql.ID # or a go string
```

This means gqlgen will be able to automatically bind to strings or ints for models you have written yourself, but the
first model in this list is used as the default type and it will always be used when:

- Generating models based on schema
- As arguments in resolvers

There isn't any way around this, gqlgen has no way to know what you want in a given context.

## Other Resources

- [Christopher Biscardi @ Gophercon UK 2018](https://youtu.be/FdURVezcdcw)
- [Introducing gqlgen: a GraphQL Server Generator for Go](https://99designs.com.au/blog/engineering/gqlgen-a-graphql-server-generator-for-go/)
- [Dive into GraphQL by Iván Corrales Solera](https://medium.com/@ivan.corrales.solera/dive-into-graphql-9bfedf22e1a)
- [Sample Project built on gqlgen with Postgres by Oleg Shalygin](https://github.com/oshalygin/gqlgen-pg-todo-example)
