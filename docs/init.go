package docs

import "strings"

func Init(tokenUrl string) {
	doc = strings.ReplaceAll(doc, "[[.TokenUrl]]", tokenUrl)
}
