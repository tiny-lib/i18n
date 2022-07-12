package pool

import (
	"sync"

	"golang.org/x/text/language"
)

var (
	LanguageTagPool = sync.Pool{New: func() interface{} {
		return language.Tag{}
	}}
)
