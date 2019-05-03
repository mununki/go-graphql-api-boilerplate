package main // import "github.com/mattdamon108/go-graphql-api-boilerplate"

import (
	"context"
	"log"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"

	"github.com/mattdamon108/go-graphql-api-boilerplate/db"
	"github.com/mattdamon108/go-graphql-api-boilerplate/handler"
	"github.com/mattdamon108/go-graphql-api-boilerplate/resolvers"
	"github.com/mattdamon108/go-graphql-api-boilerplate/schema"
)

func main() {
	db, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	context.Background()

	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}
	schema := graphql.MustParseSchema(schema.GetSchema(), &resolvers.Resolvers{DB: db}, opts...)

	mux := http.NewServeMux()
	mux.Handle("/", handler.GraphiQL{})
	mux.Handle("/query", handler.Authenticate(&handler.GraphQL{Schema: schema}))

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Listening to... port 8080")
	if err = s.ListenAndServe(); err != nil {
		panic(err)
	}
}
