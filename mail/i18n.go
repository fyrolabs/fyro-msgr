package mail

import (
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func createLocaleBundle(
	localesPath string, defaultLocale language.Tag,
) (*i18n.Bundle, error) {
	localeFiles, err := findLocaleFiles(localesPath)
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

func findLocaleFiles(localesPath string) ([]string, error) {
	pattern := filepath.Join(localesPath, "locale.*.yml")

	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	return files, nil
}
