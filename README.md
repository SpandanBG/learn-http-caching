# HTTP Caching Using HTTP Headers

## Project Setup:
```shell
go mod tidy
```

## Running the project:
```shell
make dev
```

---

## Types of Headers

### Caching Headers
There are 3 kinds of caching headers
1. expires
2. pragma
3. content-control

### Validators
1. etag
2. if-none-match
3. last-modified
4. if-modified-since

---

### 1.1 expires
Present before HTTP-1.1, the value for the `expires` header is the absolute expiry date for the data. 

Example:
```
Expires: Sat, 23 Mar 2024 12:35:58 GMT
```

The rules for this header:
1. The expiry date cannot be more than a year from the issued time.
2. The data format should be correct. (RFC1123)

> Select Caching Strategy `Expires: 5 sec from now` to view the result.

---

### 1.2 pragma
Pre HTTP-1.1 header. The only possible value for this is `no-cache`. This header is currently depricated. This was used to prevent caching. Currenly only present for backward compatibility.

Example:
```
Pragma: no-cache
```

The rules for this header:
1. The value can only be `no-cache`.
2. Is depreicated in HTTP-1.0.
3. Use `Cache-Control` instead of this.

> Select Caching Strategy `Pragma: no-cache` to view the result.

---

### 1.3 cache-control
This was introduced in HTTP-1.1. This is a multi-value/directive header.

Possible values for caching strategy:
1. `Private` - `Cache-Control: Private` will imply that the cache will be only cached in the client (browser).
2. `Public` - `Cache-Control: Public` will be cached at any of the proxies (proxy server / reverse proxy server) and is available to all the users.
3. `no-store` - `Cache-Control: no-store` will not cache anywhere.
4. `no-cache` - `Cache-Control: no-cache` will cache the data but must require validation from the server using the validation header (ETag, etc).
5. `max-age=60` - `Cache-Control: max-age=60` will cache the data for the maximum duration given in seconds. In this example for 60 seconds
6. `s-max-age=60` - `Cache-Control: s-max-age=60` same as `max-age=XYZ`, but the caching duration is for the shared places or the proxies (the `s` stands for shared).
7. `must-revalidate` - `Cache-Control: max-age=5, must-revalidate` forces the client to revalidate the cache after it is stale (based on the `max-age=XYZ` directive). If the client is not able to due to network issues, it will not serve the cache as it might be stale.
8. `proxy-revalidate` - `Cache-Control: max-age=5, must-revalidate` same as `must-revalidate` but only applies for proxies/shared-caching.

> Note: When both `max-age` and `s-max-age` is present, it will take the first for client side caching and the latter for the proxies caching

#### 1.3.1 validation headers
Along with the caching strategy that can be provided the validation strategies. These tell the server and the client how to validate the data if they are stale and/or even to use the stale data since there was no changes (the server responds with 304).

Possible values:
1. `ETag` - `ETag: "xyz"` is the entity tag header sent by the server in the response. This is used by the browser to send a `if-none-match` request header so as to pick up the cache again if the server responds with `304` and client will renew the last `max-age=N` sent by the server and start picking up from cache for that N seconds till it is stale again.
2. `Weak ETag` - `ETag: w/"xyz"` is the same as `ETag` but the value is prefixed with `w/`. What it means that the resource may be be the same, but can be regarded as the same. Ideally used by the server to either perform a strict byte by byte checking for strong ETags and doing a high-level check for weak ETags.
3. `Last-Modified` - `Last-Modified: Sun, 24 Mar 2024 13:25:24 GMT` is sent by the server in the response header. The browser then sends a `if-modified-since` request header when the cache is stale. If the server then responds with `304`, the browser then uses the cache and re-initializes the the `max-age=N` originally sent to start picking up the direcly form the cache for the next N seconds.

> Available strategies to test: (All are private caching since, I didn't create any proxies for the shared caching)

    1. Cache Control - 5 seconds
    2. Cache Control - no-store
    3. Cache Control - no-cache | ETag (controlled by query param for change of data, thus forcing revalidation logic at backend to return new value)
    4. Cache Control - 5 sec must revalidate | ETag (controlled by the same as 3, but only goes for revalidation after it is stale.)
    5. Cache Control - 5 sec must revalidate | Last-Modified (server will keep changing the data at 10 second interval)
>  All revalidate are done using HTTP status code 200/304.
