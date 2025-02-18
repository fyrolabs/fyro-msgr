package msgr

import (
	"maps"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Mail composition
type MailContents struct {
	Subject  string
	HTMLBody string
	TextBody string
}

type ComposeMailOpts struct {
	Message Message
	Locale  string
	Data    MessageData
}

func (msgr *Messenger) ComposeMail(opts ComposeMailOpts) (*MailContents, error) {
	// Merge layout data with message data
	data := maps.Clone(msgr.LayoutData)
	maps.Copy(data, opts.Data)

	// Subject
	localizer := opts.Message.Localizer(opts.Locale)
	subject, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "mail_subject",
		TemplateData: data,
	})
	if err != nil {
		return nil, err
	}

	// Body
	htmlBody := ""
	textBody := ""

	htmlLayoutTmplFile := msgr.LayoutFile(MailChannel, RenderKindHTML)
	htmlMsgTmplFiles := opts.Message.TemplateFiles(MailChannel, RenderKindHTML)
	htmlTmplFiles := append([]string{htmlLayoutTmplFile}, htmlMsgTmplFiles...)

	htmlBody, err = RenderHTML(RenderOpts{
		Templates:     htmlTmplFiles,
		Data:          data,
		Locale:        opts.Locale,
		LayoutBundle:  msgr.layoutBundle,
		MessageBundle: opts.Message.localeBundle,
	})
	if err != nil {
		return nil, err
	}

	textLayoutTmplFile := msgr.LayoutFile(MailChannel, RenderKindText)
	textMsgTmplFiles := opts.Message.TemplateFiles(MailChannel, RenderKindText)
	textTmplFiles := append([]string{textLayoutTmplFile}, textMsgTmplFiles...)

	textBody, err = RenderText(RenderOpts{
		Templates:     textTmplFiles,
		Data:          data,
		Locale:        opts.Locale,
		LayoutBundle:  msgr.layoutBundle,
		MessageBundle: opts.Message.localeBundle,
	})
	if err != nil {
		return nil, err
	}

	return &MailContents{
		Subject: subject, HTMLBody: htmlBody, TextBody: textBody,
	}, nil
}

type SMSContents struct {
	Body string
}

type ComposeSMSOpts struct {
	Message Message
	Locale  string
	Data    MessageData
}

// SMS composition
func (msgr *Messenger) ComposeSMS(opts ComposeSMSOpts) (*SMSContents, error) {
	// Merge layout data with message data
	data := maps.Clone(msgr.LayoutData)
	maps.Copy(data, opts.Data)

	// Body
	layoutTmplFile := msgr.LayoutFile(SMSChannel, RenderKindText)
	msgTmplFiles := opts.Message.TemplateFiles(MailChannel, RenderKindText)

	tmplFiles := append([]string{layoutTmplFile}, msgTmplFiles...)

	body, err := RenderText(RenderOpts{
		Templates:     tmplFiles,
		Data:          data,
		Locale:        opts.Locale,
		LayoutBundle:  msgr.layoutBundle,
		MessageBundle: opts.Message.localeBundle,
	})
	if err != nil {
		return nil, err
	}

	return &SMSContents{Body: body}, nil
}

// Push composition
type PushContents struct {
	Title string
	Body  string
}

type ComposePushOpts struct {
	Message Message
	Locale  string
	Data    MessageData
}

func (msgr *Messenger) ComposePush(opts ComposePushOpts) (*PushContents, error) {
	// Merge layout data with message data
	data := maps.Clone(msgr.LayoutData)
	maps.Copy(data, opts.Data)

	// Title
	localizer := opts.Message.Localizer(opts.Locale)
	title, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "push_title",
		TemplateData: data,
	})
	if err != nil {
		return nil, err
	}

	// Body
	layoutTmplFile := msgr.LayoutFile(SMSChannel, RenderKindText)
	msgTmplFiles := opts.Message.TemplateFiles(MailChannel, RenderKindText)

	tmplFiles := append([]string{layoutTmplFile}, msgTmplFiles...)

	body, err := RenderText(RenderOpts{
		Templates:     tmplFiles,
		Data:          data,
		Locale:        opts.Locale,
		LayoutBundle:  msgr.layoutBundle,
		MessageBundle: opts.Message.localeBundle,
	})
	if err != nil {
		return nil, err
	}

	return &PushContents{Title: title, Body: body}, nil
}
