package helper

import "golang.org/x/text/language"

func GetLangFromName(lang string) language.Tag {
	tag, err := language.Parse(lang)
	if err != nil {
		tag = language.English
	}
	return tag
}
