[![CI](https://github.com/cesbit/tlsproxy/workflows/CI/badge.svg)](https://github.com/cesbit/tlsproxy/actions)
[![Release Version](https://img.shields.io/github/release/cesbit/tlsproxy)](https://github.com/cesbit/tlsproxy/releases)

# TLS Proxy

## Installation

Just clone the project and make a build

```
git clone git@github.com:cesbit/tlsproxy.git
cd tlsproxy
go build
```

## Example usage

The following assumes a server.crt and server.key and will forward 443->80 and 8000->8000 to just-a-host.local

```
TLSPROXY_TARGET=just-a-host.local \
TLSPROXY_PORTS=443:80,8000 \
TLSPROXY_CERT_FILE=server.crt \
TLSPROXY_KEY_FILE=server.key \
tlsproxy
```

## Environment variable

Environment             | Description
----------------------- | -----------
`TLSPROXY_TARGET`       | Address of the host.
`TLSPROXY_PORTS`        | Specify the ports you want to secure with TLS. You can list multiple ports separated by commas. Use the following syntax: `<outside>:<inside>` _(example `443:80`)_. If the outside and inside ports are the same, you can simply specify the port number _(example `8000`)_.
`TLSPROXY_CERT_FILE`    | Path to the certificate file _(example `/certs/server.crt`)_.
`TLSPROXY_KEY_FILE`     | Path to the key file _(example `/certs/server.key`)_.

## Certificates

For testing, one can create certificates using the following commands:

Generate private key (.key)
Key considerations for algorithm "RSA" ≥ 2048-bit
```
openssl genrsa -out server.key 2048
```
Key considerations for algorithm "ECDSA" (X25519 || ≥ secp384r1)
https://safecurves.cr.yp.to/
List ECDSA the supported curves (openssl ecparam -list_curves)
```
openssl ecparam -genkey -name secp384r1 -out server.key
```
Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)
```
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```


