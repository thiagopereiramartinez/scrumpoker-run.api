package di

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"github.com/golobby/container"
)

func SetupFirestore() error {

	container.Singleton(func() (*firestore.Client, error) {
		ctx := context.Background()
		conf := &firebase.Config{}
		app, err := firebase.NewApp(ctx, conf)
		if err != nil {
			return nil, err
		}

		client, err := app.Firestore(ctx)
		if err != nil {
			return nil, err
		}

		return client, nil
	})

	return nil
}
