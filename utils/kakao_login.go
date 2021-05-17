package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type KakaoTokenPayload struct {
	AccessToken string `json:"access_token"`
}

type KakaoProfile struct {
	Nickname *string `json:"nickname"`
}

type KakaoAccount struct {
	Profile KakaoProfile `json:"profile"`
}

type KakaoUserInfo struct {
	Id           uint `json:"id"`
	KakaoAccount `json:"kakao_account"`
}

func KakaoLogin(code string) (*KakaoUserInfo, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	data := url.Values{
		"grant_type":   {"authorization_code"},
		"client_id":    {os.Getenv("KAKAO_REST_API_KEY")},
		"redirect_uri": {os.Getenv("KAKAO_REDIRECT_URI")},
		"code":         {code},
	}

	respToken, err := http.PostForm("https://kauth.kakao.com/oauth/token", data)
	if err != nil {
		return nil, err
	}
	defer respToken.Body.Close()

	tokenPayload := KakaoTokenPayload{}
	err = json.NewDecoder(respToken.Body).Decode(&tokenPayload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", "https://kapi.kakao.com/v2/user/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+tokenPayload.AccessToken)

	client := &http.Client{}
	respUserInfo, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer respUserInfo.Body.Close()

	userInfo := KakaoUserInfo{}
	err = json.NewDecoder(respUserInfo.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}
