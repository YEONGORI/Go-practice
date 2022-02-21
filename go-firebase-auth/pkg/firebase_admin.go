package pkg

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"

	"google.golang.org/api/option"
)

func firebaseAdmin() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("go_src/go-firebase-auth/serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, opt)

	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

}
