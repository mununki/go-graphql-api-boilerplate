package main // import "github.com/mattdamon108/go-graphql-api-boilerplate"

import (
	// "log"
	"net/http"
	// "fmt"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"

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

	defer db.DB.Close()

	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}
	schema := graphql.MustParseSchema(schema.GetSchema(), &resolvers.Resolvers{DB: db}, opts...)

	mux := http.NewServeMux()
	mux.Handle("/playground", handler.GraphiQL{})
	mux.Handle("/", &relay.Handler{Schema: schema})

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err = s.ListenAndServe(); err != nil {
		panic(err)
	}
}
