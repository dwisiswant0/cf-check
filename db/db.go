package db

import "strings"

func init() {
	Prefs = strings.Split(prefs, "\n")
}
