package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"time"
)

var ErrResolveDOHFqdn = errors.New("cannot resolve DNS over HTTPS fqdn")

type server struct {
	dns.Server
}

func newServer(ctx context.Context, provider Provider) (*server, error) {
	FQDN := provider.url.Host + "."

	query := new(dns.Msg)
	query.SetQuestion(FQDN, dns.TypeA)

	dnsClient := new(dns.Client)
	hostIPv4, _, errIPv4 := dnsClient.Exchange(query, provider.plainDNS)
	if errIPv4 != nil {
		log.Printf("cannot obtain IPv4 address for %s: %s\n", FQDN, errIPv4)
	} else {
		log.Printf("resolved %s to IPv4 %s", FQDN, hostIPv4.Answer[0].(*dns.A).A)
	}

	query.SetQuestion(FQDN, dns.TypeAAAA)
	hostIPv6, _, errIPv6 := dnsClient.Exchange(query, provider.plainDNS)
	if errIPv6 != nil {
		log.Printf("cannot obtain IPv6 address for %s: %s\n", FQDN, errIPv6)
	} else {
		log.Printf("resolved %s to IPv6 %s", FQDN, hostIPv6.Answer[0].(*dns.AAAA).AAAA)
	}

	if errIPv4 != nil && errIPv6 != nil {
		return nil, fmt.Errorf("%w: %s", ErrResolveDOHFqdn, FQDN)
	}

	server := &server{
		dns.Server{
			Addr:    ":53",
			Net:     "udp",
			Handler: newDNSHandler(ctx, 5*time.Second, provider, hostIPv4, hostIPv6),
		},
	}

	return server, nil
}
