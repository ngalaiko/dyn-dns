package main

import (
	"context"
	"flag"
	"log"
	"time"

	updater "github.com/ngalayko/dyn-dns/app"
	"github.com/ngalayko/dyn-dns/app/fetcher/ipify"
	provider "github.com/ngalayko/dyn-dns/app/provider/mock"
)

var (
	domain   = flag.String("domain", "example.com", "record domain to update")
	interval = flag.Duration("interval", time.Minute, "interval between checks")
)

func main() {
	flag.Parse()

	ctx := context.Background()

	if err := updater.New(
		provider.New(),
		ipify.New(),
		*domain,
		*interval,
	).Run(ctx); err != nil {
		log.Panicf(`[PANIC] msg="%s"`, err)
	}
}
