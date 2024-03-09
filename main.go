package main

import (
	"fmt"
	"net/http"

	"github.com/comerc/gqlgen-pg-todo/dataloaders"
	database "github.com/comerc/gqlgen-pg-todo/db"
	"github.com/comerc/gqlgen-pg-todo/graph/generated"
	"github.com/comerc/gqlgen-pg-todo/resolvers"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/fatih/color"
	"github.com/go-chi/chi"
	"github.com/go-pg/pg/v9"
)

const (
	port = ":8080"
)

func lineSeparator() {
	fmt.Println("========")
}

func startMessage() {
	lineSeparator()
	color.Green("Listening on localhost%s\n", port)
	color.Green("Visit `http://localhost%s/` in your browser\n", port)
	lineSeparator()
}

func main() {
	lineSeparator()
	// Create the database `todos` manually within postgres
	db := pg.Connect(&pg.Options{
		Database: "todos",
	})
	defer db.Close()

	fmt.Printf("%#v\n", db.Options())

	err := database.Seed(db)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	// The base path that users would use is POST /graphql which is fairly
	// idiomatic.
	r.Route("/graphql", func(r chi.Router) {
		// Initialize the dataloaders as middleware into our route
		r.Use(dataloaders.NewMiddleware(db)...)

		schema := generated.NewExecutableSchema(generated.Config{
			Resolvers: &resolvers.Resolver{
				DB: db,
			},
			Directives: generated.DirectiveRoot{},
			Complexity: generated.ComplexityRoot{},
		})

		srv := handler.NewDefaultServer(schema)
		srv.Use(extension.FixedComplexityLimit(300))

		r.Handle("/", srv)
	})

	gqlPlayground := playground.Handler("api-gateway", "/graphql")
	r.Get("/", gqlPlayground)

	startMessage()
	panic(http.ListenAndServe(port, r))
}
