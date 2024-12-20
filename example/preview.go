package test_test

import "github.com/fyrolabs/fyro-mailer/mail"

func previewLetter() {
	mailer, _ := mail.NewMailer(mail.NewMailerOpts{
		TemplatesDir:  "./templates",
		DefaultLocale: "en",
		DefaultFrom:   "sender@example.org",
		Provider: &mail.PostmarkProvider{
			ServerToken: "AAAAA",
			TrackOpens:  true,
		},
	})

	helloLetter := mail.RegisterLetterOpts{
		Name:           "helloWorld",
		ExtraTemplates: []string{},
	}

	if err := mailer.RegisterLetter(helloLetter); err != nil {
		panic(err)
	}

	data := mail.MailData{
		"Name": "John Doe",
	}

	if err := mailer.Preview(mail.PreviewOpts{
		LetterName: "helloWorld",
		Locale:     "en", // or zh-cn
		Data:       data,
	}); err != nil {
		panic(err)
	}
}
