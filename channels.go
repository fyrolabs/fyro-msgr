package msgr

type Channel string

const (
	MailChannel Channel = "mail"
	SMSChannel  Channel = "sms"
	PushChannel Channel = "push"
)

type MailChannelOpts struct {
	From    string
	ReplyTo string
}
