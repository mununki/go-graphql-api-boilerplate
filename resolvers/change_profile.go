package resolvers

import (
	"context"

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
	user := model.User{}

	if err := r.DB.First(&user, userID).Error; err != nil {
		msg := "Not existing user"
		return &ChangeProfileResponse{Status: false, Msg: &msg, User: nil}, nil
	}

	if args.Bio != nil {
		user.Bio = *args.Bio
	}
	if args.Avatar != nil {
		user.Avatar = *args.Avatar
	}

	r.DB.Save(&user)
	return &ChangeProfileResponse{Status: true, Msg: nil, User: &UserResponse{u: &user}}, nil
}

type changeProfileMutationArgs struct {
	Bio    *string
	Avatar *string
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
