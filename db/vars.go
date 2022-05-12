package db

import _ "embed"

var (
	//go:embed prefixes.txt
	prefs string
	Prefs []string
)
