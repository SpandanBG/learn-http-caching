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

1.1 expires
===
Present before HTTP-1.1, the value for the `expires` header is the absolute expiry date for the data. 

Example:
```
Expires: Sat, 23 Mar 2024 12:35:58 GMT
```

The rules for this header:
1. The expiry date cannot be more than a year from the issued time.
2. The data format should be correct. (RFC1123)

> Select Caching Strategy `Expires: 5 sec from now` to view the result.

1.2 pragma
===
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

1.3 cache-control
===
This was introduced in HTTP-1.1. This is a multi-value/directive header.

Possible values:
1. `Private` - `Cache-Control: Private` will imply that the cache will be only cached in the client (browser).
2. `Public` - `Cache-Control: Public` will be cached at any of the proxies (proxy server / reverse proxy server) and is available to all the users.
3. `no-store` - `Cache-Control: no-store` will not cache anywhere.
4. `no-cache` - `Cache-Control: no-cache` will cache the data but must require validation from the server using the validation header (ETag, etc).
5. `max-age=60` - `Cache-Control: max-age=60` will cache the data for the maximum duration given in seconds. In this example for 60 seconds
6. `s-max-age=60` - `Cache-Control: s-max-age=60` same as `max-age=XYZ`, but the caching duration is for the shared places or the proxies (the `s` stands for shared).
7. `must-revalidate` - `Cache-Control: max-age=5, must-revalidate` forces the client to revalidate the cache after it is stale (based on the `max-age=XYZ` directive). If the client is not able to due to network issues, it will not serve the cache as it might be stale.
8. `proxy-revalidate` - `Cache-Control: max-age=5, must-revalidate` same as `must-revalidate` but only applies for proxies/shared-caching.

> Note: When both `max-age` and `s-max-age` is present, it will take the first for client side caching and the latter for the proxies caching

> Available strategies to test: (All are private caching since, I didn't create any proxies for the shared caching)
    1. Cache Control - 5 seconds
    2. Cache Control - no-store
    3. Cache Control - no-cache (controlled by query param for change of data, thus forcing revalidation logic at backend to return new value)
    4. Cache Control - 5 sec must revalidate (controlled by the same as 3, but only goes for revalidation after it is stale.)
  All revalidate are done using HTTP status code 200/304.
