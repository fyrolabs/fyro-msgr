package msgr

import (
	htmlTemplate "html/template"
	"path/filepath"
	textTemplate "text/template"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func createLocaleBundle(
	path string, defaultLocale language.Tag,
) (*i18n.Bundle, error) {
	localeFiles, err := findLocaleFiles(path)
	if err != nil {
		return nil, err
	}

	bundle := i18n.NewBundle(defaultLocale)
	bundle.RegisterUnmarshalFunc("yml", yaml.Unmarshal)

	for _, file := range localeFiles {
		if _, err := bundle.LoadMessageFile(file); err != nil {
			return nil, err
		}
	}

	return bundle, nil
}

func findLocaleFiles(path string) ([]string, error) {
	pattern := filepath.Join(path, "locale.*.yml")

	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func textTemplateHelpers(
	layoutBundle *i18n.Bundle, messageBundle *i18n.Bundle, locale string,
) textTemplate.FuncMap {
	funcs := textTemplate.FuncMap{
		"tl": func(key string, data any) string {
			localizer := i18n.NewLocalizer(layoutBundle, locale)
			return localizer.MustLocalize(&i18n.LocalizeConfig{
				MessageID:    key,
				TemplateData: data,
			})
		},
		"t": func(key string, data any) string {
			localizer := i18n.NewLocalizer(messageBundle, locale)
			return localizer.MustLocalize(&i18n.LocalizeConfig{
				MessageID:    key,
				TemplateData: data,
			})
		},
	}
	return funcs
}

func htmlTemplateHelpers(
	layoutBundle *i18n.Bundle, messageBundle *i18n.Bundle, locale string,
) htmlTemplate.FuncMap {
	funcs := htmlTemplate.FuncMap{
		"tl": func(key string, data any) string {
			localizer := i18n.NewLocalizer(layoutBundle, locale)
			return localizer.MustLocalize(&i18n.LocalizeConfig{
				MessageID:    key,
				TemplateData: data,
			})
		},
		"t": func(key string, data any) string {
			localizer := i18n.NewLocalizer(messageBundle, locale)
			return localizer.MustLocalize(&i18n.LocalizeConfig{
				MessageID:    key,
				TemplateData: data,
			})
		},
	}
	return funcs
}
