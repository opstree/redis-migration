# Base image for building binary
FROM golang:1.15 as builder
LABEL VERSION=v0.1 \
      ARCH=AMD64 \
      DESCRIPTION="A redis migration tool" \
      MAINTAINER="OpsTree Solutions"
WORKDIR /go/src/redis-migrator
COPY go.mod go.mod
RUN go mod download
COPY . /go/src/redis-migrator
RUN GIT_COMMIT=$(git rev-parse --short HEAD) && \
    VERSION=$(cat ./VERSION) && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -ldflags "-X main.GitCommit=${GIT_COMMIT} -X main.Version=${VERSION}" -o redis-migrator

# Copying binary to distroless
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /go/src/redis-migrator/redis-migrator .
USER nonroot:nonroot
ENTRYPOINT ["/redis-migrator"]
