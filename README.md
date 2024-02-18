# HTTProxy Cache

Example config found [here](config.yaml):
```yaml
upstream: "https://httpbin.org"
routes:
- path: "/anything/:param1"
  clear_cache:
  - DELETE
  - POST
```

## How it works

The service has a catch-all route that proxies ALL requests to the specified upstream.  
`GET` routes defined in the `config.yamL` are cached in memory. The first request is made against upstream and all future ones are served out of memory.  
A `HEAD` request against the route will clear the cache.  
Request methods defined in `clear_cache` will also clear it.  

Example:
1. `GET /anything/foobar`
  - proxied and saved in memory
1. `GET /anything/foobar`
  - served from cache
1. `POST /anything/foobar`
  - proxied
  - clears cache
1. `GET /anything/foobar`
  - proxied and saved in memory
