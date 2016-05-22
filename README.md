Board [![Build Status](https://travis-ci.org/Lesterpig/board.svg?branch=master)](https://travis-ci.org/Lesterpig/board)
=======================================================================================================================

> One dashboard to check them all.

This repository contains a small web server used to provide a very accurate status *for all your systems*. It also supports live alerts when a service goes down.

With just a glance, you'll be able to spot the faulty parts of your infrastructure.

![Screenshot](screenshot.png "Screenshot")

Configuration
-------------

Requires a golang workspace.

```bash
cd $GOPATH/src
git clone https://github.com/lesterpig/board
cd board
cp config.sample config.go
```

At this point, you'll want to customize the `config.go` file with your very own services.

Then, build and run it!

```bash
go get ./...
board -p 8080
```

The dashboard will be available via a web interface.

Deployment
----------

You can push only the `board` binary and the `static` folder to a remote server, golang is not required for production.

Available probes
----------------

#### HTTP(S)

```go
// NewHTTP returns a ready-to-go probe.
// A warning will be triggered if the response takes more than `warning` to come.
// The `regex` is used to check the content of the website, and can be empty.
func NewHTTP(addrport string, warning time.Duration, fatal time.Duration, regex string) *HTTP
```

#### DNS

```go
// NewDNS returns a ready-to-go probe.
// `domain` will be resolved through a lookup for an A record.
// `expected` should be the first returned IPv4 address or empty to accept any IP address.
// A warning will be triggered if the response takes more than `warning` to come.
func NewDNS(addr, domain, expected string, warning, fatal time.Duration) *DNS
```

#### SMTP (over TLS)

```go
// NewSMTP returns a ready-to-go probe.
// A warning will be triggered if the response takes more than `warning` to come.
// BEWARE! Only full TLS servers are working with this probe.
func NewSMTP(addrport string, warning, fatal time.Duration) *SMTP
```

#### Minecraft

```go
// NewMinecraft returns a ready-to-go probe.
// The resulting probe will perform a real minecraft handshake to get
// some stats on the server (connected players and version).
func NewMinecraft(addrport string, fatal time.Duration) *Minecraft
```

#### Open port (TCP/UDP)

```go
// NewPort returns a ready-to-go probe.
// The `network` variable should be `tcp` or `udp` or their v4/v6 variants.
// A warning will be triggered if the response takes more than `warning` to come.
func NewPort(network, addrport string, warning, fatal time.Duration) *Port
```

Available alerts
----------------

#### Pushbullet

```go
// NewPushbullet returns a Pushbullet alerter from the private token
// available in the `account` page.
func NewPushbullet(token string) *Pushbullet
```

TODO
----

- Add probes
  + IMAP
  + OpenVPN
  + Ping (ICMP)
  + SNMP
  + DHCP
- Add alerts
  + Mail
  + Twitter
- Configure interval
- Per-service alerts
- Per-service interval
- Check more often if down
