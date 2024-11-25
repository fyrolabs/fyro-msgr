package mail

import (
	"context"

	"github.com/mrz1836/postmark"
)

type Provider interface {
	Send(opts ProviderSendOpts) error
}

type ProviderSendOpts struct {
	To       string
	From     string
	ReplyTo  string
	Subject  string
	HTMLBody string
	TextBody string
}

type PostmarkProvider struct {
	ServerToken string
	TrackOpens  bool
}

func (p *PostmarkProvider) Send(opts ProviderSendOpts) error {
	email := postmark.Email{
		From:       opts.From,
		To:         opts.To,
		ReplyTo:    opts.ReplyTo,
		Subject:    opts.Subject,
		HTMLBody:   opts.HTMLBody,
		TextBody:   opts.TextBody,
		TrackOpens: p.TrackOpens,
	}

	client := postmark.NewClient(p.ServerToken, "")
	_, err := client.SendEmail(context.Background(), email)
	return err
}
