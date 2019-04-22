package resolvers

import (
	"context"

	"github.com/mattdamon108/go-graphql-api-boilerplate/handler"
	"github.com/mattdamon108/go-graphql-api-boilerplate/model"
)

// GetMyProfile resolver
func (r *Resolvers) GetMyProfile(ctx context.Context) (*GetMyProfileResponse, error) {
	userID := ctx.Value(handler.ContextKey("userID"))

	if userID == nil {
		msg := "Not Authorized"
		return &GetMyProfileResponse{Status: false, Msg: &msg, User: nil}, nil
	}

	user := model.User{}
	if err := r.DB.DB.First(&user, userID).Error; err != nil {
		msg := "Not found"
		return &GetMyProfileResponse{Status: false, Msg: &msg, User: nil}, nil
	}
	return &GetMyProfileResponse{Status: true, Msg: nil, User: &UserResponse{u: &user}}, nil
}

// GetMyProfileResponse is the response type
type GetMyProfileResponse struct {
	Status bool
	Msg    *string
	User   *UserResponse
}

// Ok for GetMyProfileResponse
func (r *GetMyProfileResponse) Ok() bool {
	return r.Status
}

// Error for GetMyProfileResponse
func (r *GetMyProfileResponse) Error() *string {
	return r.Msg
}
