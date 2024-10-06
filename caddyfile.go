package caddyguard

import (
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(Guard{})
	httpcaddyfile.RegisterHandlerDirective("guard", parseCaddyfile)

	// order "guard" before "reverse_proxy" dir in Caddyfile
	httpcaddyfile.RegisterDirectiveOrder(PLUGIN_NAME, "before", "reverse_proxy")
}

// Parse caddy file tokens
func parseCaddyfile (h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var g Guard
	err := (&g).UnmarshalCaddyfile(h.Dispenser)

	return &g, err
}

// Parses Caddyfile syntax of Guard module.
// It's not used for real initialization, but to parse Caddy tokens and to output JSON that will be used to load the module.
// Aka some fields that are not heavy on load. Others will be loaded in Provision method.
//
//  Caddyfile-unmarshaled values will not be used directly; they will be encoded as JSON and then used from that
func (g *Guard) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {		
		if d.CountRemainingArgs() > 0 {
			return d.Err("found more arguments than it is allowed")
		}

		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "rotating_proxy":
				if !d.AllArgs(&g.Rotating_Proxy) {
					return d.Err("cannot provide more args to key 'rotating_proxy'")
				}
			case "timeout":
				timeout := ""
				if !d.AllArgs(&timeout) {
					return d.Err("cannot provide more args to key 'timeout'")
				}
				
				tD, err := time.ParseDuration(timeout)
				if err != nil {
					return err
				}

				g.Timeout = tD
			case "ip_headers":
				for d.NextArg() {
					g.IPHeaders = append(g.IPHeaders, d.Val())
				}

				for nest := d.Nesting(); d.NextBlock(nest); {
					g.IPHeaders = append(g.IPHeaders, d.Val())
				}
			case "exclude":
			case "pass_thru": 
				g.PassThrough = true

				if d.CountRemainingArgs() == 0 { continue }
				
				if d.CountRemainingArgs() > 0 {
					return d.Err("do not provide more args to 'pass_thru'")
				}

				if !d.NextArg() {
					return d.Err("cannot open block for 'pass_thru', it's standalone")
				}
			default:
				return d.Err("unknown sub-directive(s) were provided")
			}
		}
	}

	return nil
}