package provider

// Mail channel types
type MailSendOpts struct {
	To       string
	From     string
	ReplyTo  string
	Subject  string
	HTMLBody string
	TextBody string
}

type MailProvider interface {
	Send(opts MailSendOpts) error
}
