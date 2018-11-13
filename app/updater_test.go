package updater_test

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	updater "github.com/ngalayko/dyn-dns/app"
	fetcher "github.com/ngalayko/dyn-dns/app/fetcher/mock"
	providerMock "github.com/ngalayko/dyn-dns/app/provider/mock"
)

func Test_Run__should_create_new_record(t *testing.T) {
	dnsMock := providerMock.New()
	d := "test.record"
	ip := net.IPv4(127, 0, 0, 1)
	u := updater.New(
		dnsMock,
		&fetcher.Mock{IP: ip},
		d,
		time.Millisecond,
	)

	ctx := context.Background()
	defer ctx.Done()

	go func() {
		if err := u.Run(ctx); err != nil {
			t.Fatalf("can't start the app: %s", err)
		}
	}()

	time.Sleep(10 * time.Millisecond)

	domains, err := dnsMock.Get()
	if err != nil {
		t.Fatalf("can't get records: %s", err)
	}

	if len(domains) != 1 {
		log.Fatal("unexpected number of domains")
	}

	if domains[0].Name != d {
		log.Fatalf("unexpected domain name: %s", domains[0].Name)
	}

	if domains[0].Value != ip.String() {
		log.Fatalf("unexpected domain ip: %s", domains[0].Value)
	}
}

func Test_Run__should_update_existing_record(t *testing.T) {
	dnsMock := providerMock.New()
	d := "test.record"
	ip := net.IPv4(127, 0, 0, 1)

	ipMock := &fetcher.Mock{IP: ip}
	u := updater.New(
		dnsMock,
		ipMock,
		d,
		time.Millisecond,
	)

	ctx := context.Background()
	defer ctx.Done()

	go func() {
		if err := u.Run(ctx); err != nil {
			t.Fatalf("can't start the app: %s", err)
		}
	}()

	time.Sleep(10 * time.Millisecond)
	ipMock.IP = net.IPv4(1, 1, 1, 1)
	time.Sleep(10 * time.Millisecond)

	domains, err := dnsMock.Get()
	if err != nil {
		t.Fatalf("can't get records: %s", err)
	}

	if len(domains) != 1 {
		log.Fatal("unexpected number of domains")
	}

	if domains[0].Name != d {
		log.Fatalf("unexpected domain name: %s", domains[0].Name)
	}

	if domains[0].Value != ipMock.IP.String() {
		log.Fatalf("unexpected domain ip: %s", domains[0].Value)
	}
}

func Test_Run__should_handle_public_ip_err(t *testing.T) {
	dnsMock := providerMock.New()
	d := "test.record"
	var ip net.IP

	ipMock := &fetcher.Mock{IP: ip}
	u := updater.New(
		dnsMock,
		ipMock,
		d,
		time.Millisecond,
	)

	ctx := context.Background()
	defer ctx.Done()

	go func() {
		if err := u.Run(ctx); err != nil {
			t.Fatalf("can't start the app: %s", err)
		}
	}()

	time.Sleep(10 * time.Millisecond)
	ipMock.IP = net.IPv4(1, 1, 1, 1)
	time.Sleep(10 * time.Millisecond)

	domains, err := dnsMock.Get()
	if err != nil {
		t.Fatalf("can't get records: %s", err)
	}

	if len(domains) != 1 {
		log.Fatal("unexpected number of domains")
	}

	if domains[0].Name != d {
		log.Fatalf("unexpected domain name: %s", domains[0].Name)
	}

	if domains[0].Value != ipMock.IP.String() {
		log.Fatalf("unexpected domain ip: %s", domains[0].Value)
	}
}
