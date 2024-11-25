package mail

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/vanng822/go-premailer/premailer"
)

var (
	ErrInvalidKind = errors.New("invalid kind")
	ErrNoProvider  = errors.New("no send provider")
)

type MailKind struct {
	Templates []string
	Subject   map[string]string
	Data      map[string]any
}

type Mailer struct {
	Kinds        map[string]MailKind
	TemplatesDir string
	DefaultFrom  string
	SendProvider Provider
	Locale       string
}

type NewMailerOpts struct {
	TemplatesDir string
	DefaultFrom  string
	SendProvider Provider
}

func NewMailer(opts NewMailerOpts) *Mailer {
	return &Mailer{
		Kinds:        map[string]MailKind{},
		TemplatesDir: opts.TemplatesDir,
		DefaultFrom:  opts.DefaultFrom,
		SendProvider: opts.SendProvider,
		Locale:       "en", // Default locale
	}
}

func (mailer *Mailer) RegisterKind(name string, kind MailKind) {
	mailer.Kinds[name] = kind
}

type SendOpts struct {
	KindName string
	Data     map[string]any `json:"data"`
	To       string
	From     string
	ReplyTo  string
	Subject  string
	Locale   string
}

type MailData map[string]any

func (mailer *Mailer) Send(opts *SendOpts) error {
	kind, err := mailer.getKind(opts.KindName)
	if err != nil {
		return err
	}

	if mailer.SendProvider == nil {
		return ErrNoProvider
	}

	locale := mailer.Locale
	if opts.Locale != "" {
		locale = opts.Locale
	}

	htmlContent, err := mailer.BuildHTMLMessage(*kind, opts.Data, locale)
	if err != nil {
		return err
	}

	from := opts.From
	if from == "" {
		from = mailer.DefaultFrom
	}

	subject := opts.Subject
	if subject == "" {
		subject = kind.Subject[locale]
	}

	providerSendOpts := ProviderSendOpts{
		To:       opts.To,
		From:     from,
		ReplyTo:  opts.ReplyTo,
		Subject:  subject,
		HTMLBody: htmlContent,
	}

	return mailer.SendProvider.Send(providerSendOpts)
}

func (mailer *Mailer) Preview(kindName string, data MailData, locale string) error {
	kind, err := mailer.getKind(kindName)
	if err != nil {
		return err
	}

	htmlContent, err := mailer.BuildHTMLMessage(*kind, data, locale)
	if err != nil {
		return err
	}

	filePath := filepath.Join(
		"previews", fmt.Sprintf("%s_%s.html", kindName, locale),
	)
	if err := os.WriteFile(filePath, []byte(htmlContent), 0644); err != nil {
		return err
	}

	return nil
}

func (mailer *Mailer) getKind(name string) (*MailKind, error) {
	kind, exists := mailer.Kinds[name]
	if !exists {
		return nil, ErrInvalidKind
	}

	return &kind, nil
}

func (mailer *Mailer) BuildHTMLMessage(
	kind MailKind, data map[string]any, locale string,
) (string, error) {
	templates := []string{}

	for _, tmplName := range kind.Templates {
		useLocaleTmpl := true

		tmplFile := filepath.Join(
			mailer.TemplatesDir, locale, tmplName+".html.tmpl",
		)

		_, err := os.Stat(tmplFile)
		if err != nil {
			if !os.IsNotExist(err) {
				return "", err
			}
			useLocaleTmpl = false
		}

		if !useLocaleTmpl {
			tmplFile = filepath.Join(
				mailer.TemplatesDir, tmplName+".html.tmpl",
			)
		}

		templates = append(templates, tmplFile)
	}

	tmpl := template.Must(template.ParseFiles(templates...))

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {
		return "", err
	}

	premailerOpts := premailer.NewOptions()
	premailerOpts.RemoveClasses = true
	prem, err := premailer.NewPremailerFromString(buffer.String(), premailerOpts)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}
