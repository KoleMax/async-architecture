package tasks

import "strings"

func Init(tokenUrl string) {
	doc = strings.ReplaceAll(doc, "[[.TokenUrl]]", tokenUrl)
}
