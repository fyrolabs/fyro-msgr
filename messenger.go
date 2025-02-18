package msgr

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/fyrolabs/fyro-msgr/provider"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Messenger struct {
	LayoutData    MessageData
	messageMap    map[string]Message
	templatesRoot string
	mailProvider  provider.MailProvider
	mailOpts      *MailChannelOpts
	smsProvider   provider.SMSProvider
	pushProviders *provider.PushProviders
	defaultLocale language.Tag
	layoutBundle  *i18n.Bundle
}

type ClientOpts struct {
	// Path to email layout, locales, and templates
	TemplatesRoot string
	// Set the mail provider and default opts
	MailProvider provider.MailProvider
	MailOpts     *MailChannelOpts
	// SMS options
	SMSProvider provider.SMSProvider
	// Push options
	PushProviders *provider.PushProviders
	DefaultLocale string
	// Fixed data to be used in the layout
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
		mailProvider:  opts.MailProvider,
		mailOpts:      opts.MailOpts,
		smsProvider:   opts.SMSProvider,
		defaultLocale: lang, // Default locale
		LayoutData:    layoutData,
		layoutBundle:  bundle,
	}, nil
}

type AddMessageOpts struct {
	Name            string          // Must be unique
	MailChannelOpts MailChannelOpts // Email channel options
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

func (msgr *Messenger) GetMessage(name string) (*Message, error) {
	msg, exists := msgr.messageMap[name]
	if !exists {
		return nil, ErrInvalidMessage
	}

	return &msg, nil
}

type SendOpts struct {
	MessageName string
	MailTo      string                // If MailTo is defined, it will send email
	SMSTo       string                // If SMSTo is defined, it will send SMS
	PushTo      []provider.PushDevice // If pushTo has devices, it will send via push
	Data        MessageData
	Locale      string
}

func (msgr *Messenger) Send(opts SendOpts) error {
	msg, err := msgr.GetMessage(opts.MessageName)
	if err != nil {
		return err
	}

	locale := opts.Locale
	if opts.Locale == "" {
		locale = msgr.defaultLocale.String()
	}

	// errors list for each send channel
	var errs []error

	sendMail := func() error {
		from := msgr.mailOpts.From
		if msg.mailChannelOpts.From != "" {
			from = msg.mailChannelOpts.From
		}

		// Use default replyTo, unless message has its own
		replyTo := msgr.mailOpts.ReplyTo
		if msg.mailChannelOpts.ReplyTo != "" {
			replyTo = msg.mailChannelOpts.ReplyTo
		}

		contents, err := msgr.ComposeMail(ComposeMailOpts{
			Message: *msg,
			Locale:  opts.Locale,
			Data:    opts.Data,
		})
		if err != nil {
			return err
		}

		providerOpts := provider.MailSendOpts{
			To:       opts.MailTo,
			From:     from,
			ReplyTo:  replyTo,
			Subject:  contents.Subject,
			HTMLBody: contents.HTMLBody,
		}

		return msgr.mailProvider.Send(providerOpts)
	}

	sendSMS := func() error {
		contents, err := msgr.ComposeSMS(ComposeSMSOpts{
			Message: *msg,
			Locale:  locale,
			Data:    opts.Data,
		})
		if err != nil {
			return err
		}

		providerOpts := provider.SMSProviderSendOpts{
			To:   opts.SMSTo,
			Body: contents.Body,
		}

		return msgr.smsProvider.Send(providerOpts)
	}

	sendPush := func() error {
		contents, err := msgr.ComposePush(ComposePushOpts{
			Message: *msg,
			Locale:  locale,
			Data:    opts.Data,
		})
		if err != nil {
			return err
		}

		return msgr.pushProviders.Send(provider.PushProviderSendOpts{
			Devices: opts.PushTo,
			Title:   contents.Title,
			Body:    contents.Body,
		})
	}

	// Send via email
	if opts.MailTo != "" {
		if err := sendMail(); err != nil {
			errs = append(errs, err)
		}
	}

	// Send via SMS
	if opts.SMSTo != "" {
		if err := sendSMS(); err != nil {
			errs = append(errs, err)
		}
	}

	// Send via push
	if opts.PushTo != nil && msgr.pushProviders != nil {
		if err := sendPush(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (msgr *Messenger) LayoutFile(channel Channel, format RenderFormat) string {
	layout := filepath.Join(
		msgr.templatesRoot,
		fmt.Sprintf("layout_%s.%s.tmpl", channel, format),
	)
	return layout
}
