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
> Keep in mind that you need to order manually.

**If** ``pass_thru`` is provided, then there are some important headers your web server should consume:

#### X-Guard-* Headers
  - ``X-Guard-Success`` 
    > If it is set to ``1``, it means success otherwise ``-1`` means false.
  - ``X-Guard-Info``
    > Contains explainatory description.
  - ``X-Guard-Query``
    > The IP which got queried. Not present when ``X-Guard-Rate`` is ``UNKNOWN``.
  - ``X-Guard-Rate`` 
    > Either ``DANGER | LEGIT | UNKNOWN``
    > 
    > **DANGER**<br>
    > Reports that the IP reputation is bad
    >
    > **LEGIT**<br>
    > Reports that the IP reputation is good
    >
    > **LEGIT**<br>
    > Reports that the IP reputation is unknown, aka scan failure. Typically exceeded ``timeout`` constraints.


### Additional notes
Guard uses **InternetDB** to perform scans. It's completely free, and allows high traffic throughput. You can always use ``rotating_proxy`` sub-directive with Guard to increase that when needed.

Determination of a bad IP happens in the following way:
 - If **InternetDB** knows anything about the queried IP, then it is an IP with bad reputation.

### Example
```caddy
{
	order guard before reverse_proxy
}

:2000 {
	guard /api {
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

### Credits
--- Programmed by z3ntl3, will be used at the revamped ``api.pix4.dev``.
