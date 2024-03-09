package dataloaders

import (
	"context"
	"fmt"
	"net/http"

	"github.com/comerc/gqlgen-pg-todo/models"

	"github.com/go-pg/pg/v9"
	"github.com/vikstrous/dataloadgen"
)

// func User(db *pg.DB, w http.ResponseWriter, r *http.Request, next http.Handler) {
// 	loader := generated.NewUserLoader(generated.UserLoaderConfig{
// 		MaxBatch: 100,
// 		Wait:     1 * time.Millisecond,
// 		Fetch: func(keys []int) ([]*models.User, []error) {

// 			var dbUsers []*models.User
// 			// This query does NOT return an array that matches the order of the IN
// 			// clause.  Meaning: SELECT * FROM Users where id IN (1,8,3)
// 			// will not return users 1, 8, 3 in that order.  This order is VERY important
// 			// as that is how the dataloaden library resolves and matches objects.  Note the ids here are collected via
// 			// goroutines and the order is not going to be nicely ordered to match your DB
// 			// result query.  Try adding a breakpoint here and looking at the arg(keys) and the
// 			// resulting array from the following where query.
// 			err := db.Model(&dbUsers).WhereIn("id IN (?)", keys).Select()

// 			if err != nil {
// 				return []*models.User{}, []error{err}
// 			}

// 			// All we're doing here on out is just ordering our
// 			// collection to match the argument keys []int collection
// 			userKeys := make(map[int]*models.User)
// 			users := make([]*models.User, len(keys))

// 			for _, user := range dbUsers {
// 				userKeys[user.ID] = user

// 			}

// 			for i, k := range keys {
// 				if user, ok := userKeys[k]; ok {
// 					users[i] = user
// 				}
// 			}

// 			return users, []error{err}
// 		},
// 	})

// 	ctx := context.WithValue(r.Context(), UserLoader, loader)
// 	next.ServeHTTP(w, r.WithContext(ctx))
// }

func User(db *pg.DB, w http.ResponseWriter, r *http.Request, next http.Handler) {
	fetchFn := func(ctx context.Context, keys []int) ([]*models.User, []error) {
		var dbUsers []*models.User
		// This query does NOT return an array that matches the order of the IN
		// clause.  Meaning: SELECT * FROM Users where id IN (1,8,3)
		// will not return users 1, 8, 3 in that order.  This order is VERY important
		// as that is how the dataloaden library resolves and matches objects.  Note the ids here are collected via
		// goroutines and the order is not going to be nicely ordered to match your DB
		// result query.  Try adding a breakpoint here and looking at the arg(keys) and the
		// resulting array from the following where query.
		err := db.Model(&dbUsers).WhereIn("id IN (?)", keys).Select()

		if err != nil {
			return []*models.User{}, []error{err}
		}

		// All we're doing here on out is just ordering our
		// collection to match the argument keys []int collection
		userKeys := make(map[int]*models.User)
		users := make([]*models.User, len(keys))

		for _, user := range dbUsers {
			userKeys[user.ID] = user

		}

		var errs []error
		for i, k := range keys {
			if user, ok := userKeys[k]; ok {
				users[i] = user
			} else {
				// Handle the case where a key does not have a corresponding user.
				errs = append(errs, fmt.Errorf("no user found for key: %d", k))
				users[i] = nil // Keep place with nil if user is not found.
			}
		}

		return users, errs
	}

	loader := dataloadgen.NewLoader(fetchFn)

	ctx := context.WithValue(r.Context(), UserLoader, loader)
	next.ServeHTTP(w, r.WithContext(ctx))
}
