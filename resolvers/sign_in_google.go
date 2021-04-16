package resolvers

import (
	"strconv"
	"time"

	"github.com/doug-martin/goqu/v9"

	"github.com/mattdamon108/go-graphql-api-boilerplate/model"
	"github.com/mattdamon108/go-graphql-api-boilerplate/utils"
)

// SignInGoogle mutation
func (r *Resolvers) SignInGoogle(args signInGoogleMutationArgs) (*SignInGoogleResponse, error) {
	// check if exist 1) by googleID 2) by email
	// 1) googleId exists -> sign in
	// 2) email exists -> update googleId -> sign in
	// 3) otherwise sign up

	// TODO: MUST verify the google accessToken with google Auth API

	userAndSocial := model.UserAndSocial{}

	found, err := r.DB.
		From("user").
		Join(goqu.T("user_social"), goqu.On(goqu.I("user.id").Eq(goqu.I("user_social.user_id")))).
		Where(goqu.I("user_social.google").Eq(args.GoogleId)).
		ScanStruct(&userAndSocial)
	if err != nil {
		// fatal error
		msg := "Failed to sign in with google"
		return &SignInGoogleResponse{Status: false, Msg: &msg, Token: nil}, nil
	}
	if !found {
		// no matching googleId -> query by email
		found, err := r.DB.
			From("user").
			Join(goqu.T("user_social"), goqu.On(goqu.I("user.id").Eq(goqu.I("user_social.user_id")))).
			Where(goqu.I("user.email").Eq(args.Email)).
			ScanStruct(&userAndSocial)
		if err != nil {
			// fatal error
			msg := "Failed to sign in with google"
			return &SignInGoogleResponse{Status: false, Msg: &msg, Token: nil}, nil
		}
		if !found {
			// 일치하는 googleId, email 없는 경우 -> 회원가입
			tx, err := r.DB.Begin()
			if err != nil {
				msg := "Failed to sign in with google: transaction failed to begin"
				return &SignInGoogleResponse{Status: false, Msg: &msg, Token: nil}, nil
			}

			insert := tx.Insert("user").Rows(goqu.Record{"email": args.Email}).Executor()
			result, err := insert.Exec()
			if err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					msg := "Failed to sign in with google: transaction failed to be rolled back"
					return &SignInGoogleResponse{Status: false, Msg: &msg, Token: nil}, nil
				}
				msg := "Failed to sign in with google: transaction rolled back"
				return &SignInGoogleResponse{Status: false, Msg: &msg, Token: nil}, nil
			}

			id, err := result.LastInsertId()
			if err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					msg := "Failed to sign in with google: transaction failed to be rolled back"
					return &SignInGoogleResponse{Status: false, Msg: &msg, Token: nil}, nil
				}
				msg := "Failed to sign in with google: failed to insert user"
				return &SignInGoogleResponse{Status: false, Msg: &msg, Token: nil}, nil
			}

			insert = tx.Insert("user_social").Rows(goqu.Record{"user_id": id, "google": args.GoogleId}).Executor()
			if _, err := insert.Exec(); err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					msg := "Failed to sign in with google: transaction failed to be rolled back"
					return &SignInGoogleResponse{Status: false, Msg: &msg, Token: nil}, nil
				}
				msg := "Failed to sign in with kakao: transaction rolled back"
				return &SignInGoogleResponse{Status: false, Msg: &msg, Token: nil}, nil
			}

			if err = tx.Commit(); err != nil {
				msg := "Failed to sign in with kakao: transaction failed to commit"
				return &SignInGoogleResponse{Status: false, Msg: &msg, Token: nil}, nil
			}

			userIDString := strconv.Itoa(int(id))
			tokenString, err := utils.SignJWT(&userIDString)
			if err != nil {
				msg := "Error in generating JWT"
				return &SignInGoogleResponse{Status: false, Msg: &msg, Token: nil}, nil
			}

			return &SignInGoogleResponse{Status: false, Msg: nil, Token: tokenString}, nil
		} else {
			// no matching googleId, but email -> update googleId -> sign in
			insert := r.DB.
				Update("user_social").
				Set(goqu.Record{"google": args.GoogleId, "updated_at": time.Now()}).
				Where(goqu.C("user_id").Eq(userAndSocial.User.ID)).
				Executor()
			if _, err := insert.Exec(); err != nil {
				msg := "Failed to sign in with google: failed to save user social"
				return &SignInGoogleResponse{Status: false, Msg: &msg, Token: nil}, nil
			}

			userIDString := strconv.Itoa(int(userAndSocial.User.ID))
			tokenString, err := utils.SignJWT(&userIDString)
			if err != nil {
				msg := "Error in generating JWT"
				return &SignInGoogleResponse{Status: false, Msg: &msg, Token: nil}, nil
			}

			return &SignInGoogleResponse{Status: true, Msg: nil, Token: tokenString}, nil
		}
	}
	// matching googleId -> happy!
	userIDString := strconv.Itoa(int(userAndSocial.User.ID))
	tokenString, err := utils.SignJWT(&userIDString)
	if err != nil {
		msg := "Error in generating JWT"
		return &SignInGoogleResponse{Status: false, Msg: &msg, Token: nil}, nil
	}

	return &SignInGoogleResponse{Status: true, Msg: nil, Token: tokenString}, nil
}

type signInGoogleMutationArgs struct {
	GoogleId string
	Email    string
}

// SignInResponse is the response type
type SignInGoogleResponse struct {
	Status bool
	Msg    *string
	Token  *string
}

// Ok for SignUpResponse
func (r *SignInGoogleResponse) Ok() bool {
	return r.Status
}

// Error for SignUpResponse
func (r *SignInGoogleResponse) Error() *string {
	return r.Msg
}
