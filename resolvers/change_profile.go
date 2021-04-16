package resolvers

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/mattdamon108/go-graphql-api-boilerplate/handler"
	"github.com/mattdamon108/go-graphql-api-boilerplate/model"
)

// ChangeProfile mutation change profile
func (r *Resolvers) ChangeProfile(ctx context.Context, args changeProfileMutationArgs) (*ChangeProfileResponse, error) {
	userID := ctx.Value(handler.ContextKey("userID"))

	if userID == nil {
		msg := "Not Authorized"
		return &ChangeProfileResponse{Status: false, Msg: &msg, User: nil}, nil
	}

	update := r.DB.
		Update("user").
		Where(goqu.C("id").Eq(userID)).
		Set(goqu.Record{"nickname": args.Nickname, "updated_at": time.Now()}).
		Executor()
	if _, err := update.Exec(); err != nil {
		msg := "Failed to save to DB"
		return &ChangeProfileResponse{Status: false, Msg: &msg, User: nil}, nil
	}

	user := model.User{}
	found, err := r.DB.
		From("user").
		Where(goqu.C("id").Eq(userID)).
		ScanStruct(&user)
	if err != nil {
		msg := "Failed to query user"
		return &ChangeProfileResponse{Status: false, Msg: &msg, User: nil}, nil
	}
	if !found {
		msg := "Not existing user"
		return &ChangeProfileResponse{Status: false, Msg: &msg, User: nil}, nil
	}

	return &ChangeProfileResponse{Status: true, Msg: nil, User: &UserResponse{u: &user}}, nil
}

type changeProfileMutationArgs struct {
	Nickname string
}

// ChangeProfileResponse is the response type
type ChangeProfileResponse struct {
	Status bool
	Msg    *string
	User   *UserResponse
}

// Ok for ChangeProfileResponse
func (r *ChangeProfileResponse) Ok() bool {
	return r.Status
}

// Error for ChangeProfileResponse
func (r *ChangeProfileResponse) Error() *string {
	return r.Msg
}
