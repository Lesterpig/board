FROM golang:1.20.5-alpine AS builder

# Dependencies
WORKDIR /build
RUN go install github.com/GeertJohan/go.rice/rice@v1.0.3
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . ./
RUN CGO_ENABLED=0 go build -ldflags '-w -s' -o /board
RUN rice append --exec /board

# Build runtime container
FROM scratch
COPY --chown=1000:1000 --from=builder /board /board

USER 1000:1000
EXPOSE 8080

CMD ["/board"]
