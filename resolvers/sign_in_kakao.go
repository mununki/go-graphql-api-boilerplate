package resolvers

import (
	"log"
	"strconv"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/mattdamon108/go-graphql-api-boilerplate/model"
	"github.com/mattdamon108/go-graphql-api-boilerplate/utils"
)

// SignInKakao mutation
func (r *Resolvers) SignInKakao(args signInKakaoMutationArgs) (*SignInKakaoResponse, error) {
	// check if exists 1) by kakaoId
	// 1) kakaoId exists -> sign in
	// 2) input args.Email -> matching email -> update kakaoId -> sign in
	// 3) otherwise sign up

	userInfo, err := utils.KakaoLogin(args.Code)
	if err != nil {
		msg := "Failed to get a userInfo from Kakao"
		return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
	}

	userAndSocial := model.UserAndSocial{}

	found, err := r.DB.
		From("user").
		Join(goqu.T("user_social"), goqu.On(goqu.I("user.id").Eq(goqu.I("user_social.user_id")))).
		Where(goqu.I("user_social.kakao").Eq(userInfo.Id)).
		ScanStruct(&userAndSocial)
	if err != nil {
		// fatal error
		msg := "Failed to sign in with kakao"
		return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
	}
	if !found {
		if userInfo.KakaoAccount.Profile.Nickname != nil {
			// no matching kakaoId, input args.Email -> query by email
			found, err := r.DB.
				From("user").
				Join(goqu.T("user_social"), goqu.On(goqu.I("user.id").Eq(goqu.I("user_social.user_id")))).
				Where(goqu.I("user.email").Eq(userInfo.KakaoAccount.Profile.Nickname)).
				ScanStruct(&userAndSocial)
			if err != nil {
				// fatal error
				log.Panic(err)
				msg := "Failed to sign in with kakao"
				return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
			}
			if !found {
				// no matching kakaoId, no matching email -> sign up

				tx, err := r.DB.Begin()
				if err != nil {
					msg := "Failed to sign in with kakao: transaction failed to begin"
					return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
				}

				insert := tx.Insert("user").Rows(goqu.Record{"email": userInfo.KakaoAccount.Profile.Nickname}).Executor()
				result, err := insert.Exec()
				if err != nil {
					if rErr := tx.Rollback(); rErr != nil {
						msg := "Failed to sign in with kakao: transaction failed to be rolled back"
						return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
					}
					msg := "Failed to sign in with kakao: transaction rolled back"
					return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
				}

				id, err := result.LastInsertId()
				if err != nil {
					if rErr := tx.Rollback(); rErr != nil {
						msg := "Failed to sign in with kakao: transaction failed to be rolled back"
						return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
					}
					msg := "Failed to sign in with kakao: failed to insert user"
					return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
				}

				insert = tx.Insert("user_social").Rows(goqu.Record{"user_id": id, "kakao": userInfo.Id}).Executor()
				if _, err := insert.Exec(); err != nil {
					if rErr := tx.Rollback(); rErr != nil {
						msg := "Failed to sign in with kakao: transaction failed to be rolled back"
						return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
					}
					msg := "Failed to sign in with kakao: transaction rolled back"
					return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
				}

				if err = tx.Commit(); err != nil {
					msg := "Failed to sign in with kakao: transaction failed to commit"
					return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
				}

				userIDString := strconv.Itoa(int(id))
				tokenString, err := utils.SignJWT(&userIDString)
				if err != nil {
					msg := "Error in generating JWT"
					return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
				}

				return &SignInKakaoResponse{Status: true, Msg: nil, Token: tokenString}, nil
			} else {
				// no matching kakaoId, matching email -> update kakaoId -> sign in
				insert := r.DB.
					Update("user_social").
					Set(goqu.Record{"kakao": userInfo.Id, "updated_at": time.Now()}).
					Where(goqu.C("user_id").Eq(userAndSocial.User.ID)).
					Executor()
				if _, err := insert.Exec(); err != nil {
					msg := "Failed to sign in with kakao: failed to save user social"
					return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
				}

				userIDString := strconv.Itoa(int(userAndSocial.User.ID))
				tokenString, err := utils.SignJWT(&userIDString)
				if err != nil {
					msg := "Error in generating JWT"
					return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
				}

				return &SignInKakaoResponse{Status: true, Msg: nil, Token: tokenString}, nil
			}
		} else {
			// no matching kakaoId, no input args.Email -> sign up
			tx, err := r.DB.Begin()
			if err != nil {
				msg := "Failed to sign in with kakao: transaction failed to begin"
				return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
			}

			insert := tx.Insert("user").Rows(goqu.Record{"email": "no@email.com"}).Executor()
			result, err := insert.Exec()
			if err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					msg := "Failed to sign in with kakao: transaction failed to be rolled back"
					return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
				}
				msg := "Failed to sign in with kakao: transaction rolled back"
				return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
			}

			id, err := result.LastInsertId()
			if err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					msg := "Failed to sign in with kakao: transaction failed to be rolled back"
					return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
				}
				msg := "Failed to sign in with kakao: failed to insert user"
				return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
			}

			insert = tx.Insert("user_social").Rows(goqu.Record{"user_id": id, "kakao": userInfo.Id}).Executor()
			if _, err := insert.Exec(); err != nil {
				if rErr := tx.Rollback(); rErr != nil {
					msg := "Failed to sign in with kakao: transaction failed to be rolled back"
					return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
				}
				msg := "Failed to sign in with kakao: transaction rolled back"
				return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
			}

			if err = tx.Commit(); err != nil {
				msg := "Failed to sign in with kakao: transaction failed to commit"
				return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
			}

			userIDString := strconv.Itoa(int(id))
			tokenString, err := utils.SignJWT(&userIDString)
			if err != nil {
				msg := "Error in generating JWT"
				return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
			}

			return &SignInKakaoResponse{Status: true, Msg: nil, Token: tokenString}, nil
		}
	}
	// matching kakaoId -> happy!

	userIDString := strconv.Itoa(int(userAndSocial.User.ID))
	tokenString, err := utils.SignJWT(&userIDString)
	if err != nil {
		msg := "Error in generating JWT"
		return &SignInKakaoResponse{Status: false, Msg: &msg, Token: nil}, nil
	}

	return &SignInKakaoResponse{Status: true, Msg: nil, Token: tokenString}, nil
}

type signInKakaoMutationArgs struct {
	Code string
}

// SignInResponse is the response type
type SignInKakaoResponse struct {
	Status bool
	Msg    *string
	Token  *string
}

// Ok for SignUpResponse
func (r *SignInKakaoResponse) Ok() bool {
	return r.Status
}

// Error for SignUpResponse
func (r *SignInKakaoResponse) Error() *string {
	return r.Msg
}
