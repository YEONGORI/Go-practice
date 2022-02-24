package auth

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/joho/godotenv"
)

type CognitoClient struct {
	AppClientId string
	*cip.Client
}

func Init() *CognitoClient {
	error := godotenv.Load()
	if error != nil {
		log.Fatal("Error loading .env file")
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return &CognitoClient{
		AppClientId: os.Getenv("COGNITO_APP_CLIENT_ID"),
		Client:      cip.NewFromConfig(cfg),
	}
}
