package caddyguard

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

const InternetDB = "https://internetdb.shodan.io/"

var ErrBadIP = errors.New("bad ip reputation")

// Safe guards
var (
	_ caddy.Module = (*Guard)(nil)
	_ caddyhttp.MiddlewareHandler = (*Guard)(nil)
	_ caddyfile.Unmarshaler = (*Guard)(nil)
)

// Guard is an elegant IPQS plugin for Caddy.
type Guard struct {
	Timeout time.Duration	`json:"timeout,omitempty"` // If it takes longer up until timeout, will notify the web server (only) for failure with the reason, even if pass_thru is active
	IPHeaders []string `json:"ip_headers,omitempty"` // IP headers to look into to find the real ip, usefull for CDN based websites. Like Cloudflare's ``cf-connecting-ip``
	Rotating_Proxy string `json:"rotating_proxy,omitempty"` // Tells the client to use a rotating proxy when connecting to internetdb.shodan.io
	PassThrough bool `json:"pass_thru,omitempty"` // Tells whether the guard middleware should intercept strictly or pass data to the next handler, typically your web server. It does that by manipulating request headers in the form of X-Guard-* 
	logger *zap.Logger
	*http.Client
}

// CaddyModule returns the Caddy module information.
func (Guard) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.guard",
		New: func() caddy.Module { return new(Guard) },
	}
}

// Provisioning necessary parts
func (g *Guard) Provision(ctx caddy.Context) error {
	g.logger = ctx.Logger()
	g.Client = &http.Client{}

	return g.setup_client()
}

// Guard handler
func (g *Guard) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	var lookupInHeader string

	for _, header := range g.IPHeaders {
		h := r.Header.Get(header)
		if h != "" {
			lookupInHeader = h
			break
		}
	}

	if lookupInHeader == "" {
		r.Header.Add("X-Guard-Success", "-1")
		r.Header.Add("X-Guard-Info", "IP header not found")

		return next.ServeHTTP(w,r)
	} 

	err := g.Rate(lookupInHeader)

	r.Header.Add("X-Guard-Info", "Scanned IP for reputation using InternetDB")
	r.Header.Add("X-Guard-Query", lookupInHeader)
	
	switch err {
	case ErrBadIP:
		if g.PassThrough {
			r.Header.Add("X-Guard-Success", "1")
			r.Header.Add("X-Guard-Rate", "DANGER")
			
			break // exit switch - to skip statements below
		} 

		
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(403)
		
		w.Write([]byte(fmt.Sprintf(
			"<h2>You seem to use a VPN/Proxy, please turn it off to proceed.</h2>",
		)))

		return ErrBadIP
	case nil:
		r.Header.Add("X-Guard-Success", "1")
		r.Header.Add("X-Guard-Rate", "LEGIT")
	default:
		r.Header.Add("X-Guard-Success", "-1")
		r.Header.Add("X-Guard-Rate", "UNKNOWN")

		err, netErr := err.(net.Error); if netErr && err.Timeout() {
			r.Header.Set("X-Guard-Info", "Client for InternetDB timed out")
		}

		if !netErr {
			r.Header.Set("X-Guard-Info", err.Error())
		}
	}

	return next.ServeHTTP(w,r)
}


func setup_headers(req *http.Request) {
	req.Header.Add("Cache-Control", "must-revalidate")
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1")
	req.Header.Add("Content-Type", "application/json")
}

func (g *Guard) setup_client() error {
	g.Client.Timeout = g.Timeout

	if g.Rotating_Proxy != "" {
		proxyURI, err := url.Parse(g.Rotating_Proxy)
		if err != nil {
			return err
		}

		g.Client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURI),
		}
	} 

	return nil
}

func (bot *Guard) Rate(ipaddr string) error {
	req, err :=  http.NewRequest(
		http.MethodGet, 
		InternetDB + ipaddr, 
		nil,
	)

	if err != nil {
		return err
	}

	setup_headers(req)

	res, err := bot.Client.Do(req)
	if err != nil {
		return err
	}

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if !strings.Contains(string(raw), "{\"detail\":\"No information available\"}") {
		return ErrBadIP
	}

	return nil
}
