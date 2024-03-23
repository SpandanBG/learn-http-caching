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

#### 1.1 expires
Present before HTTP-1.1, the value for the `expires` header is the absolute expiry date for the data. 

Example:
```
Expires: Sat, 23 Mar 2024 12:35:58 GMT
```

The rules for this header:
1. The expiry date cannot be more than a year from the issued time.
2. The data format should be correct. (RFC1123)

> Select Caching Strategy `Expires: 5 sec from now` to view the result.

### 1.2 pragma
Pre HTTP-1.1 header. The only possible value for this is `no-cache`. This header is currently depricated. This was used to prevent caching. Currenly only present for backward compatibility.

Example:
```
Pragma: no-cache
```

The rules for this header:
1. The value can only be `no-cache`.
2. Is depreicated in HTTP-1.0.

> Select Caching Strategy `Pragma: no-cache` to view the result.
