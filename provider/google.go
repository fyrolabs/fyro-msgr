package provider

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	option "google.golang.org/api/option"
)

type GooglePushProvider struct {
	ServiceAccountKey string // path to .json key file
}

func (p *GooglePushProvider) Send(opts PushSendOpts) error {
	keyOpt := option.WithCredentialsFile(p.ServiceAccountKey)

	app, err := firebase.NewApp(context.Background(), nil, keyOpt)
	if err != nil {
		return err
	}

	ctx := context.Background()
	client, err := app.Messaging(ctx)

	message := &messaging.Message{
		Token: opts.DeviceToken,
		Notification: &messaging.Notification{
			Title: opts.Title,
			Body:  opts.Message,
		},
	}

	_, err = client.Send(ctx, message)
	if err != nil {
		return err
	}

	return nil
}
