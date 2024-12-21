package msgr

type Channel int

const (
	MailChannel Channel = iota
	SMSChannel
)

var channelMap = map[Channel]string{
	MailChannel: "mail",
	SMSChannel:  "sms",
}

func (c Channel) String() string {
	return channelMap[c]
}

type MailChannelOpts struct {
	From    string
	ReplyTo string
}
