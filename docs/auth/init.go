package auth

import "strings"

func Init(tokenUrl string) {
	doc = strings.ReplaceAll(doc, "[[.TokenUrl]]", tokenUrl)
}
