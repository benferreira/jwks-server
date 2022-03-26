# JWKS Server

[![Test](https://github.com/benferreira/jwks-server/actions/workflows/test.yml/badge.svg)](https://github.com/benferreira/jwks-server/actions/workflows/test.yml) [![codecov](https://codecov.io/gh/benferreira/jwks-server/branch/develop/graph/badge.svg?token=1Z2RFP5OIM)](https://codecov.io/gh/benferreira/jwks-server)

Basic http server that will host a JWKS generated from a provided RSA public key.

## Contents

- [JWKS Server](#jwks-server)
  - [Contents](#contents)
  - [Usage](#usage)
    - [Routes](#routes)
    - [Debug](#debug)
    - [Pretty Logging](#pretty-logging)
    - [Test Mode](#test-mode)
  - [Docker](#docker)
    - [ARM variant](#arm-variant)
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
# {"level":"info","time":"2022-03-25T23:14:29-06:00","message":"serving localhost:8000"}
```

Example response:

```sh
curl localhost:8000/api/v1/jwks.json | jq
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

### Debug

Set `DEBUG=true` to run the server with debug logging enabled:

```sh
export DEBUG=true
./bin/jwks-server
# {"level":"debug","time":"2022-03-25T23:12:52-06:00","message":"parsed *rsa.PublicKey with remaining data: \"\""}
# {"level":"info","time":"2022-03-25T23:12:52-06:00","message":"serving localhost:8000"}
```

### Pretty Logging

Enable pretty log format with the `PRETTY_LOGGING` environment variable:

```sh
export PRETTY_LOGGING=true
# 2022-03-25T23:13:50-06:00 DBG parsed *rsa.PublicKey with remaining data: ""
# 2022-03-25T23:13:50-06:00 INF serving localhost:8000
```

### Test Mode

If you want to it out without providing your own key, set the test mode environment variable, `TEST_MODE`. It will generate a random key and serve the JWKS:

```sh
export TEST_MODE=true
./bin/jwks-server
# {"level":"debug","time":"2022-03-25T23:15:01-06:00","message":"test mode enabled, generating RSA public key"}
# {"level":"debug","time":"2022-03-25T23:15:02-06:00","message":"generated public key: -----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA8G1Ac8R2W8zmqI0jRqYD\n29EIX0nPfR1n0y+M+3H+5hBujNJpoNE4wnItKnkG4F1nLNz0BBMwdN6idoNlqNGE\nngH42Nu6YLD8aWUtbdGTVDBSOkmco8sRqN4QJrT+O+PAownDVoe4+xiCA+DYO9WE\nDgLF71U9+ZUz9GF2FnFjMLDgTwJvc51SF/PwJT7RXTrbWvGWnetBQHW59tHG7M/Q\nES5RhMUUHX6ZeQTJ+soquomnDmcqTZ8PxTOT6675SAbMPvc4yk59zmkb32H+RTWp\n2wcKDNHm/kiKfJ5VYlZBfjSRapJCCjI9unWpCP7br23tbd3Gh6/Ln9JocYPOeWG/\nIQIDAQAB\n-----END PUBLIC KEY-----\n"}
# {"level":"debug","time":"2022-03-25T23:15:02-06:00","message":"parsed *rsa.PublicKey with remaining data: \"\""}
# {"level":"info","time":"2022-03-25T23:15:02-06:00","message":"serving localhost:8000"}
```

## Docker

Run this in Docker.

```sh
make build-image

docker run -p 8000:8000 -e TEST_MODE=true ko.local/jwks-server
```

### ARM variant

Build an ARM native image.

```sh
make build-image-arm
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