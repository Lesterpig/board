Board
=====

> One dashboard to check them all.

This repository contains a small web server used to provide a very accurate status *for all your systems*.

With just a glance, you'll be able to spot the faulty parts of your infrastructure.

Configuration
-------------

Requires a golang workspace.

```bash
git clone https://github.com/lesterpig/board
cd board
cp config.sample config.go
```

At this point, you'll want to customize the `config.go` file with your very own services.

Then, build and run it!

```bash
go build .
./board -p 8080
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
