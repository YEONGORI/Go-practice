package token

import (
	"context"
	"errors"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

var (
	ErrorFailedToInitializeApp   = "failed to initialize app"
	ErrorFailedToGetAuthClient   = "failed to get auth client"
	ErrorFailedToMintCostomToken = "failed to mint costom token"
)

func GenerateToken(phoneNum, password string) (string, error) {
	opt := option.WithCredentialsFile("go_src/custom-token/main-gori-341507-firebase-adminsdk-gzxf5-74cbc1d15f.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return "", errors.New(ErrorFailedToInitializeApp)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		return "", errors.New(ErrorFailedToGetAuthClient)
	}

	uid := "some-uid"
	claims := map[string]interface{}{
		"phoneNum": phoneNum,
		"password": password,
	}

	token, err := client.CustomTokenWithClaims(context.Background(), uid, claims)
	if err != nil {
		return "", errors.New(ErrorFailedToMintCostomToken)
	}

	return token, nil
}
