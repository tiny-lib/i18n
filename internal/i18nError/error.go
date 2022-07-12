package i18nError

import "errors"

var (
	LocalesDirNotFound = errors.New("locales Dir not found")
	LocalesDirNotSet   = errors.New("locales Dir not set")
	LocaleNotFound     = errors.New("locale not found")
)
