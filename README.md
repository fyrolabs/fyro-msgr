# Fyro Messenger

## Configuration
```go
type ClientOpts struct {
	// Path to email layout, locales, and templates
	TemplatesRoot string
	DefaultFrom   string
	MailProvider  provider.MailProvider
	SMSProvider   provider.SMSProvider
	DefaultLocale string
	// Dynamic data to be used in the layout
	LayoutData MessageData
}
```

Create a client instance using:
```go
msgr.NewClient(ClientOpts{})
```

## Messages
Each message has its own templates and locales.

```go
type AddMessageOpts struct {
	Name            string           // Must be unique
	MailChannelOpts *MailChannelOpts // Email channel options
}
```

 Register messages using:

```go
mailer.AddMessage(AddMessageOpts{})
```

## Templates
File structure should be as follows:

```
templates
  layout_mail.html.tmpl // Root layout
	layout_sms.html.tmpl // Root layout
  locale.en.yml
  locale.zh-cn.yml
  [entryName]/
    index_mail.html.tmpl
		index_mail.text.tmpl
		index_sms.text.tmpl
    locale.en.yml
    locale.zh-cn.yml
```

Locale files should be named locale.[lang].yml

### Subject
Message locale files must include a mandatory subject message entry, this is templated using data passed in.

```yaml
# YAML
subject: Hello {{ .Name }}
```

## Sending
```go
type SendOpts struct {
	MessageName string
	MailTo      string // If MailTo is defined, it will send email
	SMSTo       string // If SMSTo is defined, it will send SMS
	Data        MessageData
	Locale      string
}
```

## Example

Check example/preview.go for a working example
