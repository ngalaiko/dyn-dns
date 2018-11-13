package updater

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ngalayko/dyn-dns/app/fetcher"
	"github.com/ngalayko/dyn-dns/app/provider"
)

// App Updater checks for the public ip of the server
// and updates a dns record once it changed.
type App struct {
	// public ip fetcher
	ipFetcher fetcher.Fetcher
	// dns provider api
	dnsProvider provider.Provider
	// record name to change
	recordName string
	// interval between checks
	interval time.Duration
}

// New is an Updater constructor.
func New(
	p provider.Provider,
	f fetcher.Fetcher,
	recordName string,
	interval time.Duration,
) *App {
	return &App{
		dnsProvider: p,
		ipFetcher:   f,
		recordName:  recordName,
		interval:    interval,
	}
}

// Run starts updater.
func (u *App) Run(ctx context.Context) error {
	log.Printf(
		`[INFO] msg="application started" host="%s" interval="%s"`,
		u.recordName,
		u.interval,
	)

	if err := u.update(); err != nil {
		return fmt.Errorf("initial run failed: %s", err)
	}

	ticker := time.NewTicker(u.interval)
	for {
		select {
		case <-ticker.C:
			if err := u.update(); err != nil {
				log.Printf("[ERR] %s", err)
			}

		case <-ctx.Done():
			return nil
		}
	}
}

func (u *App) update() error {
	records, err := u.dnsProvider.Get()
	if err != nil {
		return err
	}

	currentIP, err := u.ipFetcher.Fetch()
	if err != nil {
		return err
	}

	for _, r := range records {
		if r.Name != u.recordName {
			continue
		}
		if r.Value == currentIP.String() {
			return nil
		}

		r.Value = currentIP.String()
		if err := u.dnsProvider.Update(r); err != nil {
			return fmt.Errorf(
				"can't update a record value to %s: %s",
				currentIP,
				err,
			)
		}
		log.Printf(
			`[INFO] msg="record updated" host="%s" ip="%s"`,
			u.recordName,
			currentIP.String(),
		)
		return nil
	}

	r := &provider.Record{
		Type:  provider.RecordTypeA,
		Name:  u.recordName,
		Value: currentIP.String(),
	}

	if err := u.dnsProvider.Create(r); err != nil {
		return fmt.Errorf(
			"can't create a record %+v: %s",
			r,
			err,
		)
	}
	log.Printf(
		`[INFO] msg="record created" host="%s" ip="%s"`,
		u.recordName,
		currentIP.String(),
	)

	return nil
}
