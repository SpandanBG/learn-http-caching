# HTTP Caching Using HTTP Headers

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

#### expires
The value for the `expires` header is the absolute expiry date for the data. 

Example:
```json
Expires: Sat, 23 Mar 2024 12:35:58 GMT
```

The rules for this header:
1. The expiry date cannot be more than a year from the issued time.
2. The data format should be correct. (RFC1123)
