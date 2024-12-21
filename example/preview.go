package example

import (
	msgr "github.com/fyrolabs/fyro-mailer"
	"github.com/fyrolabs/fyro-mailer/provider"
)

func previewLetter() {
	messenger, _ := msgr.NewClient(msgr.ClientOpts{
		TemplatesRoot: "./example/templates",
		DefaultLocale: "en",
		MailProvider: &provider.PostmarkProvider{
			ServerToken: "",
			TrackOpens:  true,
		},
	})

	userWelcome := msgr.AddMessageOpts{
		Name: "userWelcome",
		MailChannelOpts: &msgr.MailChannelOpts{
			From: "noreply@example.org",
		},
	}

	if err := messenger.AddMessage(userWelcome); err != nil {
		panic(err)
	}

	data := msgr.MessageData{
		"Name": "Bob Marley",
	}

	if err := messenger.Send(msgr.SendOpts{
		MessageName: "userWelcome",
		MailTo:      "user@example.org",
		Data:        data,
		Locale:      "en", // Locale: "en"
	}); err != nil {
		panic(err)
	}
}
