package provider

type MailProvider interface {
	Send(opts MailProviderSendOpts) error
}

type MailProviderSendOpts struct {
	To       string
	From     string
	ReplyTo  string
	Subject  string
	HTMLBody string
	TextBody string
}

type SMSProvider interface {
	Send(opts SMSProviderSendOpts) error
}

type SMSProviderSendOpts struct {
	To       string
	TextBody string
}
