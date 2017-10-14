Board [![Build Status](https://travis-ci.org/Lesterpig/board.svg?branch=master)](https://travis-ci.org/Lesterpig/board)
=======================================================================================================================

> One dashboard to check them all.

This repository contains a small web server used to provide a very accurate status *for all your systems*. It also supports live alerts when a service goes down.

With just a glance, you'll be able to spot the faulty parts of your infrastructure.

![Screenshot](screenshot.png "Screenshot")

Installation
------------

```
go get github.com/lesterpig/board
cp example.yaml board.yaml
vim board.yaml
```

Build single binary
-------------------

You may want to include the `/static/` directory as a ZIP resource in your binary.

```
go build -ldflags "-s -w" .
rice append --exec board
```

Documentation: https://github.com/GeertJohan/go.rice#append
