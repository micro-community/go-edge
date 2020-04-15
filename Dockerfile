FROM golang:1.14-alpine as builder

RUN apk --no-cache add make git gcc libtool musl-dev
WORKDIR /
COPY . /
RUN make build

FROM alpine:latest
RUN apk --no-cache add make git gcc libtool musl-dev ca-certificates dumb-init && \
    rm -rf /var/cache/apk/* /tmp/*

COPY --from=builder /x-edge .
ENTRYPOINT ["/x-edge"]
