package mail

import (
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Letter struct {
	name           string
	templatePath   string
	extraTemplates []string
	localeBundle   *i18n.Bundle
}

type NewLetterOpts struct {
	name           string
	templatePath   string
	extraTemplates []string
	defaultLocale  language.Tag
}

func NewLetter(opts NewLetterOpts) (*Letter, error) {
	bundle, err := createLocaleBundle(opts.templatePath, opts.defaultLocale)
	if err != nil {
		return nil, err
	}

	ltr := Letter{
		name:           opts.name,
		templatePath:   opts.templatePath,
		extraTemplates: opts.extraTemplates,
		localeBundle:   bundle,
	}

	return &ltr, nil
}

func (ltr *Letter) Localizer(locale string) *i18n.Localizer {
	localizer := i18n.NewLocalizer(ltr.localeBundle, locale)
	return localizer
}

func (ltr *Letter) Subject(locale string, data MailData) (string, error) {
	localizer := ltr.Localizer(locale)

	subject, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "subject",
		TemplateData: data,
	})
	if err != nil {
		return "", err
	}

	return subject, nil
}

func (ltr *Letter) TemplateFiles() []string {
	files := append([]string{"index.html.tmpl"}, ltr.extraTemplates...)
	for i, file := range files {
		files[i] = filepath.Join(ltr.templatePath, file)
	}
	return files
}
