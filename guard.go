package caddyguard

import (
	"context"
	"net/http"
	"time"

	"github.com/SimpaiX-net/ipqs"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

const InternetDB = "https://internetdb.shodan.io/"
const PLUGIN_NAME = "guard"

// Safe guards
var (
	_ caddy.Module = (*Guard)(nil)
	_ caddyhttp.MiddlewareHandler = (*Guard)(nil)
	_ caddyfile.Unmarshaler = (*Guard)(nil)
)

type mode = string

const (
	success mode = "success"
	danger mode = "danger"
	unknown mode = "unknown"
)
// Guard is an elegant IPQS plugin for Caddy.
type Guard struct { 
	TTL time.Duration 	  	`json:"ttl,omitempty"`
	Timeout time.Duration  	`json:"timeout,omitempty"`
	Proxy string 			`json:"rotating_proxy,omitempty"`
	IPHeaders []string	   	`json:"ip_headers,omitempty"`
	ctx context.Context
	logger *zap.Logger
	client *ipqs.Client
}

const ua = "Mozilla/5.0 (Linux; Android 13; SM-S901U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36"

// CaddyModule returns the Caddy module information.
func (Guard) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.guard",
		New: func() caddy.Module { return new(Guard) },
	}
}

// Provisioning necessary parts
func (g *Guard) Provision(ctx caddy.Context) error {
	ipqs.EnableCaching = true

	g.logger = ctx.Logger()
	g.client = ipqs.New().
				SetProxy(g.Proxy)

	g.ctx = context.WithValue(context.Background(), ipqs.TTL_key, g.TTL)
	return g.client.Provision()
}

// Guard handler
func (g *Guard) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) (err error) {
	defer next.ServeHTTP(w, r)

	var lookupInHeader string
	for _, header := range g.IPHeaders {
		h := r.Header.Get(header)
		if h != "" {
			lookupInHeader = h
			break
		}
	}

	if lookupInHeader == "" {
		r.Header.Set("X-Guard-Success", "-1")
		r.Header.Set("X-Guard-Info", "IP header not found")
		return 
	} 

	g.logger.Info("[GUARD-SCAN-START]:",
		zap.String("ip", lookupInHeader),
	)
	defer g.logger.Info("[GUARD-SCAN-END]:",
		zap.String("ip", lookupInHeader),
	)

	ctx, cancel := context.WithTimeout(g.ctx, g.Timeout)
	defer cancel()

	err = g.client.GetIPQS(ctx, lookupInHeader, ua)

	switch err {
		case nil:
			write_report(success, r)
		case ipqs.ErrBadIPRep:
			write_report(danger, r)
		default:
			write_report(unknown, r)
	}

	return 
}

func write_report(mode string, r *http.Request){
	switch mode {
		case success:
			r.Header.Set("X-Guard-Success", "1")
			r.Header.Set("X-Guard-Rate", "LEGIT")
		case danger:
			r.Header.Set("X-Guard-Success", "1")
			r.Header.Set("X-Guard-Rate", "DANGER")
		default:
			r.Header.Set("X-Guard-Success", "-1")
			r.Header.Set("X-Guard-Rate", "UNKNOWN")
	}
}