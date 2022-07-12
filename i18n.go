package i18n

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/czyt/i18n/internal/helper"
	"github.com/czyt/i18n/internal/i18nError"
	"github.com/czyt/i18n/internal/pool"
	"golang.org/x/exp/maps"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

const (
	localeContextKey = "locale"
)

type I18n struct {
	ctx context.Context
	// localesDir set the directory of locales
	localesDir string
	// defaultLocale is the locale default to use when no locale is specified
	defaultLocale language.Tag

	// failBackLocale set the locale when the current locale is not found
	fallBackLocale language.Tag
	printers       map[string]*message.Printer
}

func New(opt ...Opt) *I18n {
	i18n := &I18n{
		localesDir: "locales",
		printers:   make(map[string]*message.Printer),
		ctx:        context.Background(),
	}
	for _, option := range opt {
		option(i18n)
	}
	// check to avoid nil case
	return i18n
}

type Opt func(i *I18n)

func DefaultLocale(locale string) Opt {
	return func(i *I18n) {
		tag := helper.GetLangFromName(locale)
		i.defaultLocale = tag
		i.ChangeLocale(locale)
	}
}

func LocalesDir(dir string) Opt {
	return func(i *I18n) {
		i.localesDir = dir
	}
}

func DefaultFallBackLocale(locale string) Opt {
	return func(i *I18n) {
		i.fallBackLocale = helper.GetLangFromName(locale)
		catalog.Fallback(i.fallBackLocale)
	}
}

func (i *I18n) LocalesInit() error {
	builder := catalog.NewBuilder()
	// check if the locales directory is set
	if i.localesDir == "" {
		return i18nError.LocalesDirNotSet
	}
	// check if the locales directory exists
	if _, err := os.Stat(i.localesDir); os.IsNotExist(err) {
		return i18nError.LocalesDirNotFound
	}

	return filepath.Walk(i.localesDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			localeFile := info.Name()
			lang := localeFile[:len(localeFile)-len(filepath.Ext(localeFile))]
			payload := pool.PayloadPool.Get().(map[string]string)

			defer func() {
				maps.Clear(payload)
				pool.PayloadPool.Put(payload)
			}()
			data, err := ioutil.ReadFile(filepath.Join(i.localesDir, localeFile))
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, &payload)
			if err != nil {
				return err
			}
			tag, parseErr := language.Parse(lang)
			if parseErr != nil {
				tag = language.Make(lang)
			}
			for key, value := range payload {
				errSetString := builder.SetString(tag, key, value)
				if errSetString == nil {
					continue
				}
			}
			i.printers[tag.String()] = message.NewPrinter(tag, message.Catalog(builder))
		}
		return nil
	})

}

func (i *I18n) ChangeLocale(locale string) {
	i.ctx = context.WithValue(i.ctx, localeContextKey, locale)
}

// Trf translate the source string with the given parameters
func (i *I18n) Trf(source string, param ...interface{}) string {
	return i.getCurrentPrinter().Sprintf(source, param...)
}

// TrfWriter translate the source string with the given parameters and write it to the writer
func (i *I18n) TrfWriter(writer io.Writer, source string, param ...interface{}) (int, error) {
	return i.getCurrentPrinter().Fprintf(writer, source, param...)
}

// TrPrint translate and print the source string
func (i *I18n) TrPrint(source string) (int, error) {
	return i.getCurrentPrinter().Print(source)
}

// TrfPrint translate and print the source string with the given writer
func (i *I18n) TrfPrint(writer io.Writer, source string) (int, error) {
	return i.getCurrentPrinter().Fprint(writer, source)
}

// TrPrintln translate and println the source string
func (i *I18n) TrPrintln(source string) (int, error) {
	return i.getCurrentPrinter().Println(source)
}

// TrfPrintln translate and println the source string with the given writer
func (i *I18n) TrfPrintln(writer io.Writer, source string) (int, error) {
	return i.getCurrentPrinter().Fprint(writer, source)
}

func (i *I18n) getCurrentPrinter() *message.Printer {
	tag := pool.LanguageTagPool.Get().(language.Tag)
	defer func() {
		pool.LanguageTagPool.Put(tag)
	}()
	if localeName, ok := i.ctx.Value(localeContextKey).(string); ok {
		if localeName == "" {
			tag = i.defaultLocale
		} else {
			tag = helper.GetLangFromName(localeName)
		}
	} else {
		tag = i.defaultLocale
	}
	return i.printers[tag.String()]
}

func (i *I18n) GetPrinterByLocale(localeName string) (*message.Printer, error) {
	tag := helper.GetLangFromName(localeName)
	if _, ok := i.printers[tag.String()]; !ok {
		return nil, i18nError.LocaleNotFound
	}
	return i.printers[tag.String()], nil
}
