# JWKS Server

[![Test](https://github.com/benferreira/jwks-server/actions/workflows/test.yml/badge.svg)](https://github.com/benferreira/jwks-server/actions/workflows/test.yml) [![codecov](https://codecov.io/gh/benferreira/jwks-server/branch/develop/graph/badge.svg?token=1Z2RFP5OIM)](https://codecov.io/gh/benferreira/jwks-server)

Basic http server that will host a JWKS generated from a provided RSA public key.

## Contents

- [JWKS Server](#jwks-server)
  - [Contents](#contents)
  - [Usage](#usage)
    - [Routes](#routes)
    - [Flags](#flags)
    - [Test Mode](#test-mode)
  - [Build](#build)
  - [Test](#test)

## Usage

```sh
export RSA_PUB_KEY='-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0j69to/Sy6u5HeZfBdBt
esvO618W80K3HEo/mnapyhEJsqTQzLy5F+0OM1ZkZCrdFpXuN6jRNHP4tAVEdwg/
aY1uezwbHFwC80HLVNPTmshn5Q4zwFTX60lzSR43BkzwGQPNGixIbkNpLyPBJb0F
yRn9bSbNbwR8GtXYLBNCj4k5ITCb2fWezcxFSajcXfHZPKQgKkXa80P3YEgL/dFU
KFV6gQSmDBm/5nlU4KmUW8loo9AVAL91jW6N2msDlrVg6t8y0aG1w0kAQDwr/sFp
C6W+e6LzeskeuNAvbtgl3blKakSsZ/BB6JIZz+r6YJGFeU0KGYTcpyxSC+T/nJzd
uQIDAQAB
-----END PUBLIC KEY-----'

./bin/jwks-server
# 11:36AM INF serving localhost:8000
```

Example response:

```sh
> curl localhost:8000/api/v1/jwks.json | jq
# {
#   "keys": [
#     {
#       "alg": "RS256",
#       "kty": "RSA",
#       "n": "0j69to/Sy6u5HeZfBdBtesvO618W80K3HEo/mnapyhEJsqTQzLy5F+0OM1ZkZCrdFpXuN6jRNHP4tAVEdwg/aY1uezwbHFwC80HLVNPTmshn5Q4zwFTX60lzSR43BkzwGQPNGixIbkNpLyPBJb0FyRn9bSbNbwR8GtXYLBNCj4k5ITCb2fWezcxFSajcXfHZPKQgKkXa80P3YEgL/dFUKFV6gQSmDBm/5nlU4KmUW8loo9AVAL91jW6N2msDlrVg6t8y0aG1w0kAQDwr/sFpC6W+e6LzeskeuNAvbtgl3blKakSsZ/BB6JIZz+r6YJGFeU0KGYTcpyxSC+T/nJzduQ",
#       "e": "AQAB"
#     }
#   ]
# }
```

### Routes

The following routes are served:

* /health
* /api/v1/jwks.json

### Flags

```sh
./bin/jwks-server -h    
# Usage of ./bin/jwks-server:
#   -d    enable debug logging
#   -p int
#         http listening port (default 8000)
#   -t    generate an RSA key and serve a JWKS
```

### Test Mode

If you want to it out without providing your own key, set the test flag, `-t`. It will generate a random key and serve the JWKS:

```sh
./bin/jwks-server -t -d
# 11:41AM DBG test flag provided, generating RSA public key
# 11:41AM DBG generated public key: -----BEGIN PUBLIC KEY-----
# MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsJfMtANIRRTG8GOVaeDx
# KUhf6zvPI4Rwa1YDLdisPTbtI9BsfZhQGXQPYc226wgN9xJR/J7xyA/9T6y2Rykq
# TtA6dnF50GW6C8w+Pb7J+ufFEkJXABUmZQMvZQaFbAHveV4yfqdwMh8YiEVr4ZPx
# 9Idz1+NAdzOO5R8OnePikoThvPHBsmsW3Jdd5sjrno4Skb8NfoE8PnSxleQ4pGMK
# o7q/bF0jPNzaX8G0dcCVUPo2YiRhULWWMzLHSTxZab5KGX88nlWc5wcCXhnqPyKS
# WicX3mY6jEAkned3Rk7plum7DYgp+I0D/C7mxLZwG4JHry4Pq9ITnC+/bBvzu/RA
# eQIDAQAB
# -----END PUBLIC KEY-----

# 11:41AM DBG parsed *rsa.PublicKey with remaining data: ""
# 11:41AM INF serving localhost:8000
```

## Build

Requirements:
* Golang >= 1.17.x


```sh
make build
```

## Test

```sh
make test
```