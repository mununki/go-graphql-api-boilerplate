package resolvers

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"

	"github.com/mattdamon108/go-graphql-api-boilerplate/handler"
	"github.com/mattdamon108/go-graphql-api-boilerplate/utils"
)

// ChangePassword mutation change password
func (r *Resolvers) ChangePassword(ctx context.Context, args changePasswordMutationArgs) (*ChangePasswordResponse, error) {
	userID := ctx.Value(handler.ContextKey("userID"))

	if userID == nil {
		msg := "Not Authorized"
		return &ChangePasswordResponse{Status: false, Msg: &msg}, nil
	}

	hashed, err := utils.HashPassword(args.Password)
	if err != nil {
		msg := "Failed to hash the password"
		return &ChangePasswordResponse{Status: false, Msg: &msg}, nil
	}

	update := r.DB.Update("user").Where(goqu.C("id").Eq(userID)).Set(goqu.Record{"password": hashed, "updated_at": time.Now()}).Executor()
	if _, err := update.Exec(); err != nil {
		msg := "Failed to save to DB"
		return &ChangePasswordResponse{Status: false, Msg: &msg}, nil
	}

	return &ChangePasswordResponse{Status: true, Msg: nil}, nil
}

type changePasswordMutationArgs struct {
	Password string
}

// ChangePasswordResponse is the response type
type ChangePasswordResponse struct {
	Status bool
	Msg    *string
}

// Ok for ChangePasswordResponse
func (r *ChangePasswordResponse) Ok() bool {
	return r.Status
}

// Error for ChangePasswordResponse
func (r *ChangePasswordResponse) Error() *string {
	return r.Msg
}
