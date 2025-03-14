package provider

import (
	"fmt"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"github.com/sideshow/apns2/token"
)

type ApplePushProvider struct {
	KeyID             string
	TeamID            string
	PrivateKey        string
	NotificationTopic string
}

func (p *ApplePushProvider) Send(opts PushSendOpts) error {
	alertData := map[string]string{"body": opts.Message}

	if opts.Title != "" {
		// Add subtitle if it is set
		alertData["subtitle"] = opts.Title
	}

	pl := payload.NewPayload().Alert(alertData)

	authKey, err := token.AuthKeyFromFile(p.PrivateKey)
	if err != nil {
		return err
	}

	token := &token.Token{
		AuthKey: authKey,
		KeyID:   p.KeyID,
		TeamID:  p.TeamID,
	}

	client := apns2.NewTokenClient(token)

	notification := apns2.Notification{
		DeviceToken: opts.DeviceToken,
		Topic:       p.NotificationTopic,
		Payload:     pl,
	}

	res, err := client.Push(&notification)
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}
