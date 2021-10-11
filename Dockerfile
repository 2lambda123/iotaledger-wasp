ARG GOLANG_IMAGE_TAG=1.17-buster

# Build stage
FROM golang:${GOLANG_IMAGE_TAG} AS build

ARG BUILD_TAGS=rocksdb,builtin_static

RUN mkdir /wasp
WORKDIR /wasp

# Make sure that modules only get pulled when the module file has changed
COPY go.mod go.sum /wasp/
RUN go mod download
RUN go mod verify

# Project build stage
COPY . .

RUN go build -tags=${BUILD_TAGS}
RUN go build -tags=${BUILD_TAGS} ./tools/wasp-cli

# Testing stages
# Complete testing
# FROM golang:1.16.5-buster AS test-full
# WORKDIR /run

# COPY --from=build $GOPATH/pkg/mod $GOPATH/pkg/mod
# COPY --from=build /wasp/ /run

# CMD go test -tags rocksdb -timeout 20m ./...

# # Unit tests without integration tests
# FROM golang:1.16.5-buster AS test-unit
# WORKDIR /run

# COPY --from=build $GOPATH/pkg/mod $GOPATH/pkg/mod
# COPY --from=build /wasp/ /run

# CMD go test -tags rocksdb -short ./...

# Wasp CLI build
# FROM golang:1.16.5-buster as wasp-cli
# COPY --from=build /wasp/wasp-cli /usr/bin/wasp-cli
# ENTRYPOINT ["wasp-cli"]

# Wasp build
FROM gcr.io/distroless/cc
FROM golang:${GOLANG_IMAGE_TAG}

# Config is overridable via volume mount to /run/config.json
# COPY docker_config.json /run/config.json

COPY --from=build /wasp/wasp /usr/bin/wasp
COPY --from=build /wasp/wasp-cli /usr/bin/wasp-cli

# ENTRYPOINT ["/usr/bin/wasp"]
