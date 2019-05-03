package resolvers

import (
	"github.com/mattdamon108/go-graphql-api-boilerplate/model"
)

// SignUp mutation creates user
func (r *Resolvers) SignUp(args signUpMutationArgs) (*SignUpResponse, error) {

	newUser := model.User{Email: args.Email, Password: args.Password, FirstName: args.FirstName, LastName: args.LastName}

	if !r.DB.Where("email = ?", args.Email).First(&model.User{}).RecordNotFound() {
		msg := "Already signed up"
		return &SignUpResponse{Status: false, Msg: &msg, User: nil}, nil
	}

	newUser.HashPassword()

	r.DB.Create(&newUser)

	return &SignUpResponse{Status: true, Msg: nil, User: &UserResponse{u: &newUser}}, nil
}

type signUpMutationArgs struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

// SignUpResponse is the response type
type SignUpResponse struct {
	Status bool
	Msg    *string
	User   *UserResponse
}

// Ok for SignUpResponse
func (r *SignUpResponse) Ok() bool {
	return r.Status
}

// Error for SignUpResponse
func (r *SignUpResponse) Error() *string {
	return r.Msg
}
