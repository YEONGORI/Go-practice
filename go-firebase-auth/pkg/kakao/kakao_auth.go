package kakao

import "net/http"

func GetAuthCode(w http.ResponseWriter, r *http.Request) {
	redirectUrl := "https://kauth.kakao.com/oauth/authorize?client_id=6c428a4d1ebc925a93e5a2336dc31b97&redirect_uri=Whttps://l7adsepqf9.execute-api.us-east-2.amazonaws.com/beta/kakao/signin&response_type=code"
	http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
	res, _ := http.Get(redirectUrl)
	url, _ := res.Location()

	val := url.Query()
	authCode := val.Get("code")
}

func CreateCustomToken()
