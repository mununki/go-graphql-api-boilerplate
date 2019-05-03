package resolvers

import (
	"github.com/mattdamon108/go-graphql-api-boilerplate/db"
)

// Resolvers including query and mutation
type Resolvers struct {
	*db.DB
}
