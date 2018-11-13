package main

import (
	"context"
	"flag"
	"log"
	"net"
	"time"

	updater "github.com/ngalayko/dyn-dns/app"
	fetcher "github.com/ngalayko/dyn-dns/app/fetcher/mock"
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
		&fetcher.Mock{IP: net.IPv4(127, 0, 0, 1)},
		*domain,
		*interval,
	).Run(ctx); err != nil {
		log.Panic(err)
	}
}
