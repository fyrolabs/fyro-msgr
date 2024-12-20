package mail

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/exp/maps"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Mailer struct {
	letterMap     map[string]Letter
	templatesDir  string
	defaultFrom   string
	provider      Provider
	defaultLocale language.Tag
	layoutBundle  *i18n.Bundle
	layoutData    MailData
}

type NewMailerOpts struct {
	// Path to email layout, locales, and templates
	TemplatesDir  string
	DefaultFrom   string
	Provider      Provider
	DefaultLocale string
	// Dynamic data to be used in the layout
	LayoutData MailData
}

func NewMailer(opts NewMailerOpts) (*Mailer, error) {
	if opts.Provider == nil {
		return nil, ErrInvalidProvider
	}

	lang, err := language.Parse(opts.DefaultLocale)
	if err != nil {
		return nil, err
	}

	bundle, err := createLocaleBundle(opts.TemplatesDir, lang)
	if err != nil {
		return nil, err
	}

	layoutData := MailData{}
	if opts.LayoutData != nil {
		layoutData = opts.LayoutData
	}

	return &Mailer{
		letterMap:     map[string]Letter{},
		templatesDir:  opts.TemplatesDir,
		defaultFrom:   opts.DefaultFrom,
		provider:      opts.Provider,
		defaultLocale: lang, // Default locale
		layoutData:    layoutData,
		layoutBundle:  bundle,
	}, nil
}

type RegisterLetterOpts struct {
	Name           string
	ExtraTemplates []string
}

func (mailer *Mailer) RegisterLetter(opts RegisterLetterOpts) error {
	ltr, err := NewLetter(NewLetterOpts{
		name:           opts.Name,
		templatePath:   filepath.Join(mailer.templatesDir, opts.Name),
		extraTemplates: opts.ExtraTemplates,
		defaultLocale:  mailer.defaultLocale,
	})
	if err != nil {
		return err
	}

	mailer.letterMap[opts.Name] = *ltr
	return nil
}

func (mailer *Mailer) getLetter(name string) (*Letter, error) {
	ltr, exists := mailer.letterMap[name]
	if !exists {
		return nil, ErrInvalidKind
	}

	return &ltr, nil
}

type SendOpts struct {
	LetterName string
	Data       MailData `json:"data"`
	To         string
	From       string
	ReplyTo    string
	Subject    string
	Locale     string
}

type MailData map[string]any

func (mailer *Mailer) Send(opts SendOpts) error {
	ltr, err := mailer.getLetter(opts.LetterName)
	if err != nil {
		return err
	}

	locale := opts.Locale
	if opts.Locale == "" {
		locale = mailer.defaultLocale.String()
	}

	from := opts.From
	if from == "" {
		from = mailer.defaultFrom
	}

	// Compose the letter
	contents, err := mailer.Compose(ComposeOpts{
		Letter: *ltr,
		Locale: locale,
		Data:   opts.Data,
	})
	if err != nil {
		return err
	}

	subject := opts.Subject
	if opts.Subject == "" {
		subject = contents.Subject
	}

	providerSendOpts := ProviderSendOpts{
		To:       opts.To,
		From:     from,
		ReplyTo:  opts.ReplyTo,
		Subject:  subject,
		HTMLBody: contents.HTMLMessage,
	}

	return mailer.provider.Send(providerSendOpts)
}

type ComposeOpts struct {
	Letter Letter
	Locale string
	Data   MailData
}

type LetterContents struct {
	Subject     string
	HTMLMessage string
	TextMessage string
}

func (mailer *Mailer) Compose(opts ComposeOpts) (*LetterContents, error) {
	subject, err := opts.Letter.Subject(opts.Locale, opts.Data)
	if err != nil {
		return nil, err
	}

	// Merge layout data with letter data
	mergedData := maps.Clone(mailer.layoutData)
	maps.Copy(mergedData, opts.Data)

	layoutTmplFile := filepath.Join(mailer.templatesDir, "layout.html.tmpl")
	ltrTmplFiles := opts.Letter.TemplateFiles()

	tmplFiles := append([]string{layoutTmplFile}, ltrTmplFiles...)

	htmlMessage, err := RenderHTMLetter(RenderOpts{
		Templates:    tmplFiles,
		Data:         mergedData,
		Locale:       opts.Locale,
		LayoutBundle: mailer.layoutBundle,
		LetterBundle: opts.Letter.localeBundle,
	})
	if err != nil {
		return nil, err
	}

	return &LetterContents{
		Subject:     subject,
		HTMLMessage: htmlMessage,
	}, nil
}

type PreviewOpts struct {
	LetterName string
	Data       MailData
	Locale     string
}

func (mailer *Mailer) Preview(opts PreviewOpts) error {
	ltr, err := mailer.getLetter(opts.LetterName)
	if err != nil {
		return err
	}

	contents, err := mailer.Compose(ComposeOpts{
		Letter: *ltr,
		Locale: opts.Locale,
		Data:   opts.Data,
	})
	if err != nil {
		return err
	}

	filePath := filepath.Join(
		"previews", fmt.Sprintf("%s_%s.html", opts.LetterName, opts.Locale),
	)
	if err := os.WriteFile(filePath, []byte(contents.HTMLMessage), 0644); err != nil {
		return err
	}

	return nil
}
