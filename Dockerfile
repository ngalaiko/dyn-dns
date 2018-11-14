FROM golang:1.11-alpine3.7 as builder

ARG GOOS
ARG GOARCH

WORKDIR $GOPATH/src/github.com/ngalayko/dyn-dns

COPY . .

RUN GOOS=$GOOS GOARCH=$GOARCH go build -o /dyn-dns ./cmd/dyn-dns

FROM alpine:3.7

RUN apk add --no-cache --update \
        ca-certificates

COPY --from=builder /dyn-dns /dyn-dns

ENTRYPOINT [ "/dyn-dns" ]
