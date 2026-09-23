package main

import (
	"log"
	_ "net/http/pprof"

	_ "modernc.org/sqlite"

	"capuchin/internal/app/capuchin"
)

var appVersion = "v0.0.0" //nolint:gochecknoglobals // все в порядке

func main() {
	err := capuchin.Start(appVersion)
	if err != nil {
		log.Fatalf("Start: %v\n", err)
	}
}
