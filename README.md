# TCS Status Board

> One dashboard to check them all.

This repository contains a small web server used to provide a very accurate status *for all your systems*.
It also supports live alerts when a service goes down.

With just a glance, you'll be able to spot the faulty parts of your infrastructure.

![Screenshot](screenshot.png "Screenshot")

This project is forked from [`Lesterpig/board`](https://github.com/Lesterpig/board) and modified to fit the needs of the Trifork Cloud Stack.

## Installation

GoReleaser is used to build and release the binaries and Docker containers.
You can find the latest release [here](https://github.com/trifork/tcs-status/releases), and the latest container image [here](https://github.com/trifork/tcs-status/pkgs/container/tcs-status).

To run the board, you need to create a `board.yaml` file.
You can find an example in the repository in [`examples/board.yaml`](./examples/board.yaml).

This needs to be placed in the same directory as where you run the board from.
To mount the file into the container and expose the board on port 8080, you can use the following command:

```bash
docker run -v "${PWD}"/board.yaml:/board.yaml:ro -p 8080:8080 ghcr.io/trifork/tcs-board:latest
```

In [`examples/`](./examples/), you will find Kubernetes manifests to deploy the board to a Kubernetes cluster, which makes use of a ConfigMap containing the `board.yaml` file.

## Development

We have added a Makefile to simplify the development process.
You can use `make build` to build the binary, and `make run` to run the board locally.

We make use of [`rice`](https://github.com/GeertJohan/go.rice) to embed the static files in [static/](./static/) into the binary.
When running `make build`, the static files will be embedded into the built binary.
When running `make run`, the static files will be served directly from the [static/](./static/) directory.

To build the Docker container, you can use `make image`.
This will create a container image tagged with `tcs-board:latest`, which can be run using `docker run  -v "${PWD}/examples/board.yaml:/board.yaml:ro" -p 8080:8080 ghcr.io/trifork/tcs-board:0.0.1-next`.

To deploy the container to a Kubernetes cluster, you can use 

```bash
make cluster # Creates a local Kubernetes cluster using k3d (requires Docker)
make load-image # Load the container image into the local cluster
make deploy # Deploy the board to the local cluster
```

This will create a local Kubernetes cluster using [k3d](https://k3d.io), and deploy the board to it.
The board will be exposed on port  http://localhost:8081. If you want to re-redeploy the board, you can use `make redeploy`.

Any binaries used in the Makefile are installed using `go install` and can be found in `bin/`.


## Release

We use [GoReleaser](https://github.com/goreleaser/goreleaser) to release new versions of the tool in binary format to Github Releases and built to Docker containers to [`ghcr.io/trifork/tcs-board`](https://github.com/trifork/tcs-board/pkgs/container/tcs-board).

To release a new version, manually trigger the [Continuous Delivery workflow](https://github.com/trifork/tcs-board/actions/workflows/release.yaml) with a new version in SemVer format (`x.y.z`)
