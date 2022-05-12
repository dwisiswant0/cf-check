package main

import (
	"bufio"
	"sync"

	"github.com/dwisiswant0/cf-check/db"
)

var (
	wg sync.WaitGroup
	sc *bufio.Scanner

	domainMode, showCloudflare bool

	cidrs = db.Prefs
)
