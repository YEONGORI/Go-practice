package token

import (
	"context"
	"errors"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/internal"

	"google.golang.org/api/option"
)

var (
	ErrorFailedToInitializeApp   = "failed to initialize app"
	ErrorFailedToGetAuthClient   = "failed to get auth client"
	ErrorFailedToMintCostomToken = "failed to mint costom token"
)

func GenerateToken(email, password string) (string, error) {
	opt := option.WithCredentialsFile("go_src/custom-token/serviceAccount.json")
	config := &firebase.Config{ProjectID: "main-gori-341507"}
	_, err := firebase.NewApp(context.Background(), config, opt)

	if err != nil {
		return "", errors.New(ErrorFailedToInitializeApp)
	}

	configAuth := &internal.AuthConfig{
		ProjectID:        "main-gori-341507",
		ServiceAccountID: "109426012739542871894",
	}

	client, err := auth.NewClient(context.Background(), configAuth)
	// client, err := app.Auth(context.Background())
	if err != nil {
		return "", errors.New(ErrorFailedToGetAuthClient)
	}

	uid := "some-uid"
	claims := map[string]interface{}{
		"email":    email,
		"password": password,
	}

	token, err := client.CustomTokenWithClaims(context.Background(), uid, claims)
	if err != nil {
		return "", errors.New(ErrorFailedToMintCostomToken)
	}

	return token, nil
}
