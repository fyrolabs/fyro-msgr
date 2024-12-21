package msgr

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"

	"github.com/fyrolabs/fyro-mailer/provider"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Messenger struct {
	messageMap    map[string]Message
	templatesRoot string
	defaultFrom   string
	mailProvider  provider.MailProvider
	smsProvider   provider.SMSProvider
	defaultLocale language.Tag
	layoutBundle  *i18n.Bundle
	layoutData    MessageData
}

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

type MessageData map[string]any

func NewClient(opts ClientOpts) (*Messenger, error) {
	if opts.MailProvider == nil && opts.SMSProvider == nil {
		return nil, ErrNoProviders
	}

	lang, err := language.Parse(opts.DefaultLocale)
	if err != nil {
		return nil, err
	}

	bundle, err := createLocaleBundle(opts.TemplatesRoot, lang)
	if err != nil {
		return nil, err
	}

	layoutData := MessageData{}
	if opts.LayoutData != nil {
		layoutData = opts.LayoutData
	}

	return &Messenger{
		messageMap:    map[string]Message{},
		templatesRoot: opts.TemplatesRoot,
		defaultFrom:   opts.DefaultFrom,
		mailProvider:  opts.MailProvider,
		smsProvider:   opts.SMSProvider,
		defaultLocale: lang, // Default locale
		layoutData:    layoutData,
		layoutBundle:  bundle,
	}, nil
}

type AddMessageOpts struct {
	Name            string
	MailChannelOpts *MailChannelOpts
}

func (msgr *Messenger) AddMessage(opts AddMessageOpts) error {
	msg, err := NewMessage(NewMessageOpts{
		name:            opts.Name,
		templatesPath:   filepath.Join(msgr.templatesRoot, opts.Name),
		mailChannelOpts: opts.MailChannelOpts,
		defaultLocale:   msgr.defaultLocale,
	})
	if err != nil {
		return err
	}

	msgr.messageMap[opts.Name] = *msg
	return nil
}

func (msgr *Messenger) getMessage(name string) (*Message, error) {
	msg, exists := msgr.messageMap[name]
	if !exists {
		return nil, ErrInvalidKind
	}

	return &msg, nil
}

type SendOpts struct {
	MessageName string
	MailTo      string
	SMSTo       string
	Data        MessageData
	Locale      string
}

func (msgr *Messenger) Send(opts SendOpts) error {
	msg, err := msgr.getMessage(opts.MessageName)
	if err != nil {
		return err
	}

	locale := opts.Locale
	if opts.Locale == "" {
		locale = msgr.defaultLocale.String()
	}

	// Send via email
	if opts.MailTo != "" {
		subject, err := msg.MailSubject(locale, opts.Data)
		if err != nil {
			return err
		}

		content, err := msgr.Compose(ComposeOpts{
			Message: *msg,
			Channel: MailChannel,
			Format:  RenderKindHTML,
			Locale:  locale,
			Data:    opts.Data,
		})
		if err != nil {
			return err
		}

		providerOpts := provider.MailProviderSendOpts{
			To:       opts.MailTo,
			From:     msg.mailChannelOpts.From,
			ReplyTo:  msg.mailChannelOpts.ReplyTo,
			Subject:  subject,
			HTMLBody: content,
		}

		return msgr.mailProvider.Send(providerOpts)
	}

	// Send via SMS
	if opts.SMSTo != "" {
		content, err := msgr.Compose(ComposeOpts{
			Message: *msg,
			Channel: SMSChannel,
			Format:  RenderKindText,
			Locale:  locale,
			Data:    opts.Data,
		})
		if err != nil {
			return err
		}

		providerOpts := provider.SMSProviderSendOpts{
			To:       opts.SMSTo,
			TextBody: content,
		}

		return msgr.smsProvider.Send(providerOpts)
	}

	return nil
}

func (msgr *Messenger) LayoutFile(channel Channel, format RenderFormat) string {
	layout := filepath.Join(
		msgr.templatesRoot,
		fmt.Sprintf("layout_%s.%s.tmpl", channel.String(), format),
	)
	return layout
}

type ComposeOpts struct {
	Message Message
	Channel Channel
	Format  RenderFormat
	Locale  string
	Data    MessageData
}

func (msgr *Messenger) Compose(opts ComposeOpts) (string, error) {
	// Merge layout data with message data
	data := maps.Clone(msgr.layoutData)
	maps.Copy(data, opts.Data)

	layoutTmplFile := msgr.LayoutFile(opts.Channel, opts.Format)
	msgTmplFiles := opts.Message.TemplateFiles(opts.Channel, opts.Format)

	tmplFiles := append([]string{layoutTmplFile}, msgTmplFiles...)

	var content string
	var err error

	if opts.Format == RenderKindText {

	} else if opts.Format == RenderKindHTML {
		content, err = RenderHTML(RenderOpts{
			Templates:     tmplFiles,
			Data:          data,
			Locale:        opts.Locale,
			LayoutBundle:  msgr.layoutBundle,
			MessageBundle: opts.Message.localeBundle,
		})

	} else {
		return "", ErrInvalidFormat
	}

	return content, err
}

type PreviewOpts struct {
	MessageName string
	Data        MessageData
	Locale      string
}

func (msgr *Messenger) Preview(opts PreviewOpts) error {
	ltr, err := msgr.getMessage(opts.MessageName)
	if err != nil {
		return err
	}

	content, err := msgr.Compose(ComposeOpts{
		Message: *ltr,
		Locale:  opts.Locale,
		Data:    opts.Data,
	})
	if err != nil {
		return err
	}

	filePath := filepath.Join(
		"previews", fmt.Sprintf("%s_%s.html", opts.MessageName, opts.Locale),
	)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return err
	}

	return nil
}
