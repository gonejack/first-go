package main

import (
	"bytes"
	"fmt"
	"github.com/miekg/dns"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"context"
	"errors"
)

var ErrHTTPStatus = errors.New("bad HTTP status")

type handler struct {
	context        context.Context
	provider       Provider
	hostIPv4       *dns.Msg
	hostIPv6       *dns.Msg
	client         *http.Client
	httpBufferPool *sync.Pool
	udpBufferPool  *sync.Pool
}

func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	if h.askingForProviderIP(r) {
		// Answer back with IP address obtained at start
		h.replyWithProviderIP(w, r)
		return
	}

	buffer := h.udpBufferPool.Get().([]byte)
	// no need to reset buffer as wire takes care of cutting it down
	wire, err := r.PackBuffer(buffer)
	if err != nil {
		log.Printf("cannot pack message to wire format: %s\n", err)
		dns.HandleFailed(w, r)
		return
	}

	// It's fine to copy the slice headers as long as we keep
	// the underlying array of bytes.
	h.udpBufferPool.Put(buffer) //nolint:staticcheck

	respWire, err := h.requestHTTP(h.context, wire)
	if err != nil {
		log.Printf("HTTP request failed: %s\n", err)
		dns.HandleFailed(w, r)
		return
	}

	message := new(dns.Msg)
	if err := message.Unpack(respWire); err != nil {
		log.Printf("cannot unpack message from wireformat: %s\n", err)
		dns.HandleFailed(w, r)
		return
	}

	message.SetReply(r)
	if err := w.WriteMsg(message); err != nil {
		log.Println("write dns message error: ", err)
	}
}

func (h *handler) askingForProviderIP(r *dns.Msg) bool {
	return len(r.Question) > 0 && r.Question[0].Name == h.provider.url.Host+"." &&
		(r.Question[0].Qtype == dns.TypeA || r.Question[0].Qtype == dns.TypeAAAA)
}

func (h *handler) replyWithProviderIP(w dns.ResponseWriter, r *dns.Msg) {
	host := h.hostIPv4
	if r.Question[0].Qtype == dns.TypeAAAA {
		host = h.hostIPv6
	}

	host.SetReply(r)
	if err := w.WriteMsg(host); err != nil {
		log.Println("write dns message error: ", err)
	}
}

func (h *handler) requestHTTP(context context.Context, wire []byte) (respWire []byte, err error) {
	buffer := h.httpBufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer h.httpBufferPool.Put(buffer)

	_, err = buffer.Write(wire)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(context, http.MethodPost, h.provider.url.String(), buffer)
	if err != nil {
		return nil, err
	}

	const contentTypeUDPWireFormat = "application/dns-udpwireformat"
	request.Header.Set("Content-Type", contentTypeUDPWireFormat)

	response, err := h.client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %s", ErrHTTPStatus, response.Status)
	}

	respWire, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if err := response.Body.Close(); err != nil {
		return nil, err
	}

	return respWire, nil
}

func newDNSHandler(
	context context.Context,
	httpTimeout time.Duration,
	provider Provider,
	hostIPv4, hostIPv6 *dns.Msg,
) dns.Handler {
	client := &http.Client{
		Timeout: httpTimeout,
	}

	httpBufferPool := &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(nil)
		},
	}

	const udpPacketMaxSize = 512
	udpBufferPool := &sync.Pool{
		New: func() interface{} {
			return make([]byte, udpPacketMaxSize)
		},
	}

	return &handler{
		context:        context,
		provider:       provider,
		hostIPv4:       hostIPv4,
		hostIPv6:       hostIPv6,
		client:         client,
		httpBufferPool: httpBufferPool,
		udpBufferPool:  udpBufferPool,
	}
}
