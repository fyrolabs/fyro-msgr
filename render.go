package msgr

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/vanng822/go-premailer/premailer"
)

type RenderFormat string

var (
	RenderKindText RenderFormat = "text"
	RenderKindHTML RenderFormat = "html"
)

type RenderOpts struct {
	Templates     []string
	Data          MessageData
	Locale        string
	LayoutBundle  *i18n.Bundle
	MessageBundle *i18n.Bundle
}

func RenderText(opts RenderOpts) (string, error) {
	funcs := textTemplateHelpers(
		opts.LayoutBundle, opts.MessageBundle, opts.Locale,
	)

	tmplName := filepath.Base(opts.Templates[0])

	tmpl, err := template.New(tmplName).
		Funcs(funcs).ParseFiles(opts.Templates...)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, opts.Data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func RenderHTML(opts RenderOpts) (string, error) {
	funcs := htmlTemplateHelpers(
		opts.LayoutBundle, opts.MessageBundle, opts.Locale,
	)

	tmplName := filepath.Base(opts.Templates[0])

	tmpl, err := template.New(tmplName).
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
