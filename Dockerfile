FROM golang:1.20.5-alpine AS build-env

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
FROM alpine:3.18.4
WORKDIR /app
COPY --from=build-env /board /app/board
EXPOSE 8080
CMD ["/app/board"]
