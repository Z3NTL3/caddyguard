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
    <img alt="Static Badge" src="https://img.shields.io/badge/z3ntl3-white?style=flat&logo=aeromexico&logoSize=auto&label=Author">

</div>

## Intro
--- **Guard** is an elegant IPQS plugin for Caddy. Acting as a middleware or microservice between your web server.

--- **Features** are built in, you can tell Guard to intercept or pass data all the way down to your web server.

--- **Questions?** feel free to ask by [contacting me](https://t.me/z3ntl3)! 

### Install
```
xcaddy build --with github.com/z3ntl3/caddyguard
```

### Example
```caddy

:2000 {
    # guard is ordered before "reverse_proxy"
    # https://caddyserver.com/docs/caddyfile/directives#directive-order
	guard /api* {
		rotating_proxy 1.1.1.1 
		timeout 3s 
		ip_headers cf-connecting-ip {
			more1
			more2
			more3
		}
		pass_thru 
	}

	reverse_proxy  http://localhost:2000
}
```


### Caddyfile syntax
```caddy
guard [matcher] {
    rotating_proxy <arg>
    timeout <arg>
    ip_headers <args...> {
        <arg> 
        <arg>
        <arg>
        ...
    }
    pass_thru 
}
```
### Sub-directives 
 - ``rotating_proxy <arg> ``
     > **Doc**
     > - Should comfort [net.http](https://pkg.go.dev/net/http#Transport). 
     >
     > - Supported protocols are ``socks``, ``http`` and ``https``.
     > - If scheme is not provided, ``http`` is assumed.
     >
     > **Examples**
     >
     > ```caddy
     >  guard /api* {
     >     rotating_proxy 1.1.1.1    
     >  }
     > ```
     > Here ``http://1.1.1.1`` is assumed.
     >
     > <br>
     >
     > ```caddy
     >  guard /api* {
     >     rotating_proxy socks5://1.1.1.1    
     >  }
     > ```
     > Here ``socks5://1.1.1.1`` is assumed.
     >
     > **NOTE**<br>
     > Underlying client may change. [Proxifier](https://github.com/Z3NTL3/proxifier) > may be binded to this plugin. Which is our own low-level proxy client library.
     > 
 - ``timeout <arg> ``
     > **Doc**
     > - Should comfort [time](https://pkg.go.dev/time#ParseDuration). 
     > 
     > Aka arg values like: ``10s``, ``1m`` etc...
     >
     > **Examples**
     >
     > ```caddy
     >  guard /api* {
     >     timeout 1s    
     >  }
     
- ``ip_headers <args...> {...}``
    > **Doc**
     > - Can be arbitrary values. Tells Guard plugin to find the real ip address in one of those headers.
     > 
     > Values like: ``cf-connecting-ip``, ``x-forwarded-for`` and etc..., seem logical
     >
     > **Examples**
     >
     > ```caddy
     >  guard /api* {
     >     ip_headers header1 {
     >          header2
     >     }
     >  }
- ``pass_thru``
    > **Doc**
    > <br>
    > Accepts no args, and disallows opening a block. It acts like a turn on.
    >
    > Providing ``pass_thru`` simply means to pass data down to the next handler, aka your web server/reverse proxy instead of writing a builtin HTTP response. Provides useful data by manipulating the request headers, while it moves down to the next tree of handler(s).
    > 
    > **If** ``pass_thru`` is provided, then there are some important headers your web server should consume:
    >
    > #### X-Guard-* Headers
    > - ``X-Guard-Success`` 
    >     > If it is set to ``1``, it means success otherwise ``-1`` means false.
    > - ``X-Guard-Info``
    >     > Contains explainatory description.
    > - ``X-Guard-Query``
    >     > The IP which got queried. Not present when ``X-Guard-Rate`` is ``UNKNOWN``.
    > - ``X-Guard-Rate`` 
    >     > Either ``DANGER | LEGIT | UNKNOWN``
    >     > 
    >     > **DANGER**<br>
    >     > Reports that the IP reputation is bad
    >     >
    >     > **LEGIT**<br>
    >     > Reports that the IP reputation is good
    >     >
    >     > **UNKNOWN**<br>
    >     > Reports that the IP reputation is unknown, aka scan failure. Typically exceeded ``timeout`` constraints.


### Additional notes
Guard uses **InternetDB** to perform scans. It's completely free, and allows high traffic throughput. You can always use ``rotating_proxy`` sub-directive with Guard to allow a limitless quota when needed.

Determination of a bad IP happens based on logic:
 - If **InternetDB** has information regards the queried IoT device by IP, then consider it has a bad IP reputation.


### Credits
--- Programmed by [z3ntl3](https://z3ntl3.com)
