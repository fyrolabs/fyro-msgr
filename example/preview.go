package example

import (
	msgr "github.com/fyrolabs/fyro-msgr"
	"github.com/fyrolabs/fyro-msgr/preview"
	"github.com/fyrolabs/fyro-msgr/provider"
)

func Preview() {
	client, _ := msgr.NewClient(msgr.ClientOpts{
		TemplatesRoot: "./example/templates",
		DefaultLocale: "en",
		MailProvider: &provider.PostmarkProvider{
			ServerToken: "",
			TrackOpens:  true,
		},
	})

	userWelcome := msgr.AddMessageOpts{
		Name: "userWelcome",
		MailChannelOpts: msgr.MailChannelOpts{
			From: "noreply@example.org",
		},
	}

	if err := client.AddMessage(userWelcome); err != nil {
		panic(err)
	}

	data := msgr.MessageData{
		"Name": "Bob Marley",
	}

	if err := preview.PreviewMessage(client, preview.PreviewOpts{
		MessageName: "userWelcome",
		Channels: []msgr.Channel{
			msgr.MailChannel, msgr.SMSChannel, msgr.PushChannel,
		},
		Data:   data,
		Locale: "en",
		OutDir: "./example/out",
	}); err != nil {
		panic(err)
	}
}
