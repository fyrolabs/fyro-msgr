package mail

import (
	"bytes"
	"text/template"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/vanng822/go-premailer/premailer"
)

type RenderOpts struct {
	Templates    []string
	Data         MailData
	Locale       string
	LayoutBundle *i18n.Bundle
	LetterBundle *i18n.Bundle
}

func RenderHTMLetter(opts RenderOpts) (string, error) {
	funcs := template.FuncMap{
		"tl": func(key string, data any) string {
			localizer := i18n.NewLocalizer(opts.LayoutBundle, opts.Locale)
			return localizer.MustLocalize(&i18n.LocalizeConfig{
				MessageID:    key,
				TemplateData: data,
			})
		},
		"t": func(key string, data any) string {
			localizer := i18n.NewLocalizer(opts.LetterBundle, opts.Locale)
			return localizer.MustLocalize(&i18n.LocalizeConfig{
				MessageID:    key,
				TemplateData: data,
			})
		},
	}

	tmpl, err := template.New("layout.html.tmpl").
		Funcs(funcs).ParseFiles(opts.Templates...)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, opts.Data); err != nil {
		return "", err
	}

	premailerOpts := premailer.NewOptions()
	premailerOpts.RemoveClasses = true
	prem, err := premailer.NewPremailerFromString(buffer.String(), premailerOpts)
	if err != nil {
		return "", err
	}

	htmlMessage, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return htmlMessage, nil
}
