package i18n

import (
	"log"
)

func ExampleTr() {
	i18nCall := New(
		DefaultLocale("zh_Hans"),
		LocalesDir("testdata/locales"),
		DefaultFallBackLocale("en_US"),
	)
	err := i18nCall.LocalesInit()
	if err != nil {
		panic(err)
	}
	i18nCall.ChangeLocale("zh_SC")
	tr := i18nCall.Trf("Score", "张三", 88)
	log.Println(tr)
	i18nCall.ChangeLocale("zh_Hans")
	tr = i18nCall.Trf("Score", "张三", 88)
	log.Println(tr)
}
