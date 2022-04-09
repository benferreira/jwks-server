# JWKS Server

[![Test](https://github.com/benferreira/jwks-server/actions/workflows/test.yml/badge.svg)](https://github.com/benferreira/jwks-server/actions/workflows/test.yml) [![codecov](https://codecov.io/gh/benferreira/jwks-server/branch/main/graph/badge.svg?token=1Z2RFP5OIM)](https://codecov.io/gh/benferreira/jwks-server)

Basic http server that will host a JWKS generated from a provided RSA public key.

## Contents

- [JWKS Server](#jwks-server)
  - [Contents](#contents)
  - [Usage](#usage)
    - [Single public key](#single-public-key)
    - [Multiple Keys](#multiple-keys)
    - [Routes](#routes)
    - [Debug](#debug)
    - [Pretty Logging](#pretty-logging)
    - [Test Mode](#test-mode)
  - [Docker](#docker)
    - [ARM variant](#arm-variant)
  - [Build](#build)
  - [Test](#test)

## Usage

### Single public key

If you have a single public key, you can set `RSA_PUB_KEY` and the server will generate a JWKS from it.

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
```

```json
{"level":"info","time":"2022-03-25T23:14:29-06:00","message":"serving localhost:8000"}
```

Example response:

```sh
curl localhost:8000/api/v1/jwks.json
```

```json
{
  "keys": [
    {
      "alg": "RS256",
      "kty": "RSA",
      "n": "0j69to/Sy6u5HeZfBdBtesvO618W80K3HEo/mnapyhEJsqTQzLy5F+0OM1ZkZCrdFpXuN6jRNHP4tAVEdwg/aY1uezwbHFwC80HLVNPTmshn5Q4zwFTX60lzSR43BkzwGQPNGixIbkNpLyPBJb0FyRn9bSbNbwR8GtXYLBNCj4k5ITCb2fWezcxFSajcXfHZPKQgKkXa80P3YEgL/dFUKFV6gQSmDBm/5nlU4KmUW8loo9AVAL91jW6N2msDlrVg6t8y0aG1w0kAQDwr/sFpC6W+e6LzeskeuNAvbtgl3blKakSsZ/BB6JIZz+r6YJGFeU0KGYTcpyxSC+T/nJzduQ",
      "e": "AQAB"
    }
  ]
}
```

### Multiple Keys

If you have multiple keys and/or would like to provide a `kid` for each key, set the path for the keys file with `RSA_KEYS_FILE`. The yaml keys file uses the following format:

```yaml
rsaPubKeys:
  - key: |
      publicKey1
    kid: key1
  - key: |
      publicKey2
    kid: key2
```

Full example:

```sh
cat << EOF > ./keys.yaml
rsaPubKeys:
  - key: |
      -----BEGIN PUBLIC KEY-----
      MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0j69to/Sy6u5HeZfBdBt
      esvO618W80K3HEo/mnapyhEJsqTQzLy5F+0OM1ZkZCrdFpXuN6jRNHP4tAVEdwg/
      aY1uezwbHFwC80HLVNPTmshn5Q4zwFTX60lzSR43BkzwGQPNGixIbkNpLyPBJb0F
      yRn9bSbNbwR8GtXYLBNCj4k5ITCb2fWezcxFSajcXfHZPKQgKkXa80P3YEgL/dFU
      KFV6gQSmDBm/5nlU4KmUW8loo9AVAL91jW6N2msDlrVg6t8y0aG1w0kAQDwr/sFp
      C6W+e6LzeskeuNAvbtgl3blKakSsZ/BB6JIZz+r6YJGFeU0KGYTcpyxSC+T/nJzd
      uQIDAQAB
      -----END PUBLIC KEY-----
    kid: key1
  - key: |
      -----BEGIN PUBLIC KEY-----
      MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA9/YhUFEO8X7VkWEbduO8
      gVYS0sHW2G1hyQxVLPY/RPM/mbZ1az3WZTRjRpYtnjQrKn1KHJ671xzn8GfSWRgD
      rkaNF6mTzFI5BTaGORQS+8emFV2xXgVo0pagRFlBdhB2P3xkzNkpJzeKaDD+lG1q
      t1eyJcZTvurSBF9C6nytqdbIc7YIbd8dcYLt3zoCGNm+KBl6aXnpeD6AGt7C1pIj
      yvSWPaY4MUJ6beC6A70PIQxKLgECY9eWLwBqnceUO7EPo9aZuq7tDwE7h7h5TRtT
      pJZ9RTyK6UEiZ/v1e4cKVEVvLL9rmFXLTeJKXYIZbI1lYbuanipiFI/VNvv6r2lF
      CwIDAQAB
      -----END PUBLIC KEY-----
    kid: key2
EOF

export RSA_KEYS_FILE=./keys.yaml
./bin/jwks-server
```

```log
{"level":"info","time":"2022-04-09T13:53:59-06:00","message":"serving localhost:8000"}
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
```

```json
{"level":"debug","time":"2022-03-25T23:12:52-06:00","message":"parsed *rsa.PublicKey with remaining data: \"\""}
{"level":"info","time":"2022-03-25T23:12:52-06:00","message":"serving localhost:8000"}
```

### Pretty Logging

Enable pretty log format with the `PRETTY_LOGGING` environment variable:

```sh
export PRETTY_LOGGING=true
./bin/jwks-server
```

```log
2022-03-25T23:13:50-06:00 DBG parsed *rsa.PublicKey with remaining data: ""
2022-03-25T23:13:50-06:00 INF serving localhost:8000
```

### Test Mode

If you want to it out without providing your own key, set the test mode environment variable, `TEST_MODE`. It will generate a random key and serve the JWKS:

```sh
export TEST_MODE=true
./bin/jwks-server
```

```json
{"level":"debug","time":"2022-03-25T23:15:01-06:00","message":"test mode enabled, generating RSA public key"}
{"level":"debug","time":"2022-03-25T23:15:02-06:00","message":"generated public key: -----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA8G1Ac8R2W8zmqI0jRqYD\n29EIX0nPfR1n0y+M+3H+5hBujNJpoNE4wnItKnkG4F1nLNz0BBMwdN6idoNlqNGE\nngH42Nu6YLD8aWUtbdGTVDBSOkmco8sRqN4QJrT+O+PAownDVoe4+xiCA+DYO9WE\nDgLF71U9+ZUz9GF2FnFjMLDgTwJvc51SF/PwJT7RXTrbWvGWnetBQHW59tHG7M/Q\nES5RhMUUHX6ZeQTJ+soquomnDmcqTZ8PxTOT6675SAbMPvc4yk59zmkb32H+RTWp\n2wcKDNHm/kiKfJ5VYlZBfjSRapJCCjI9unWpCP7br23tbd3Gh6/Ln9JocYPOeWG/\nIQIDAQAB\n-----END PUBLIC KEY-----\n"}
{"level":"debug","time":"2022-03-25T23:15:02-06:00","message":"parsed *rsa.PublicKey with remaining data: \"\""}
{"level":"info","time":"2022-03-25T23:15:02-06:00","message":"serving localhost:8000"}
```

## Docker

Run this in Docker.

```sh
make build-image

docker run --rm -p 8000:8000 -e TEST_MODE=true ko.local/jwks-server
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