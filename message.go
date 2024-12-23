package msgr

import (
	"fmt"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Message struct {
	name            string
	templatePath    string
	mailChannelOpts MailChannelOpts
	localeBundle    *i18n.Bundle
}

type NewMessageOpts struct {
	name            string
	templatesPath   string
	mailChannelOpts MailChannelOpts
	defaultLocale   language.Tag
}

func NewMessage(opts NewMessageOpts) (*Message, error) {
	bundle, err := createLocaleBundle(opts.templatesPath, opts.defaultLocale)
	if err != nil {
		return nil, err
	}

	msg := Message{
		name:            opts.name,
		templatePath:    opts.templatesPath,
		mailChannelOpts: opts.mailChannelOpts,
		localeBundle:    bundle,
	}

	return &msg, nil
}

func (msg *Message) Localizer(locale string) *i18n.Localizer {
	localizer := i18n.NewLocalizer(msg.localeBundle, locale)
	return localizer
}

func (msg *Message) MailSubject(
	locale string, data MessageData,
) (string, error) {
	localizer := msg.Localizer(locale)

	subject, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "subject",
		TemplateData: data,
	})
	if err != nil {
		return "", err
	}

	return subject, nil
}

func (msg *Message) TemplateFiles(channel Channel, format RenderFormat) []string {
	// TODO implement partials

	index := filepath.Join(
		msg.templatePath, fmt.Sprintf("index_%s.%s.tmpl", channel.String(), format),
	)

	return []string{index}
}
