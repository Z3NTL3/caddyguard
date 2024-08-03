<!-- header -->

<div align="center">   
    <div>
        <img src="/img/logo.png" width=300 style="border: 2px solid grey;"><br>
         <div>
                <img alt="GitHub License" src="https://img.shields.io/github/license/z3ntl3/caddyguard" >
                <img alt="GitHub Repo stars" src="https://img.shields.io/github/stars/z3ntl3/caddyguard">
                <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/z3ntl3/caddyguard">
        </div>
    </div>
</div>

## Intro
--- **Guard** is an elegant plugin for Caddy. It provides IP reputation scan. Acting as a middleware between your web server.

--- **Features** are built in, you can tell Guard to intercept or to pass data to your web server.



### Caddyfile syntax
```
guard [matcher] {
    rotating_proxy 1.1.1.1 # not required
    timeout 3s
    ip_headers cf-connecting-ip {
        more1
        more2
        more3
    }
    pass_thru # not required
}
```
> You need to order manually.
