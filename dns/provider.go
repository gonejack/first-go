package main

import "net/url"

type Provider struct {
	plainDNS string
	url      url.URL
}

func Cloudflare() Provider {
	return Provider{
		plainDNS: "1.1.1.1:53",
		url: url.URL{
			Scheme: "https",
			Host:   "cloudflare-dns.com",
			Path:   "/dns-query",
		},
	}
}
