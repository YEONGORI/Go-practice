package pkg

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"

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
