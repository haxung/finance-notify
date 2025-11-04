FROM alpine:3.20 as root-certs

RUN apk add -U --no-cache ca-certificates
RUN addgroup -g 1001 hax
RUN adduser hax -u 1001 -D -G hax /home/hax

FROM golang:1.21 as builder

WORKDIR /build
COPY --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o notify

FROM scratch as final

COPY --from=root-certs /etc/passwd /etc/passwd
COPY --from=root-certs /etc/group /etc/group
COPY --chown=1001:1001 --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --chown=1001:1001 --from=builder /build/notify /notify

RUN mkdir /log && chown 1001:1001 /log

VOLUME [ "/log" ]

USER hax
ENTRYPOINT ["/notify"]
STOPSIGNAL SIGQUIT
