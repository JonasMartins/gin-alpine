package utils

import (
	"embed"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"gopkg.in/yaml.v3"
)

//go:embed locales/*.yaml
var localeFS embed.FS

type Translator struct {
	Bundle *i18n.Bundle
	Lang   string
}

func NewTranslator(defaultLang string) *Translator {
	bundle := i18n.NewBundle(language.Portuguese) // idioma padrão
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	// Carrega arquivos embutidos
	_, err := bundle.LoadMessageFileFS(localeFS, "locales/pt.yaml")
	if err != nil {
		FatalResult("unable to load locale file", err)
	}
	_, err = bundle.LoadMessageFileFS(localeFS, "locales/en.yaml")
	if err != nil {
		FatalResult("unable to load locale file", err)
	}

	return &Translator{Bundle: bundle, Lang: defaultLang}
}

func (t *Translator) T(key string, data map[string]any) string {
	localizer := i18n.NewLocalizer(t.Bundle, t.Lang)
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: data,
	})
	if err != nil {
		return key // fallback → mostra a chave
	}
	return msg
}
