//go:generate go run github.com/99designs/gqlgen --verbose
package resolvers

import (
	"github.com/comerc/gqlgen-pg-todo/graph/generated"

	"github.com/go-pg/pg/v9"
)

type Resolver struct {
	DB *pg.DB
}

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Todo() generated.TodoResolver {
	return &todoResolver{r}
}

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }
