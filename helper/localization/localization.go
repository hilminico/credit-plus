package localization

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"path/filepath"
	"strings"
)

var (
	bundle *i18n.Bundle
)

const (
	DefaultLanguage = "en"
	LangKey         = "lang"
)

func InitLocalization() error {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Load all locale files
	pattern := filepath.Join("locale", "*.toml")
	matches, err := filepath.Glob(pattern)

	if err != nil {
		return fmt.Errorf("failed to find locale files: %v", err)
	}

	for _, filename := range matches {
		if _, err := bundle.LoadMessageFile(filename); err != nil {
			return fmt.Errorf("failed to load locale file %s: %v", filename, err)
		}
	}

	return nil
}

func GetLanguage(ctx context.Context) string {
	if lang, ok := ctx.Value(LangKey).(string); ok && lang != "" {
		return normalizeLanguage(lang)
	}
	return DefaultLanguage
}

// normalizeLanguage ensures consistent language format
func normalizeLanguage(lang string) string {
	// Convert to lowercase and take first part (e.g., "en-US" -> "en")
	lang = strings.ToLower(lang)
	if parts := strings.Split(lang, "-"); len(parts) > 0 {
		return parts[0]
	}
	return lang
}

func Localize(ctx context.Context, messageID string, templateData map[string]interface{}) string {
	lang := GetLanguage(ctx)
	localizer := i18n.NewLocalizer(bundle, lang)

	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})

	if err != nil {
		return messageID // Fallback to message ID
	}

	return message
}

func WithLanguage(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, LangKey, lang)
}

func GetLocalizer(ctx context.Context) *i18n.Localizer {
	acceptLanguage := GetLanguage(ctx)
	return i18n.NewLocalizer(bundle, acceptLanguage)
}
