FROM golang:1.21-alpine3.18 as builder

RUN go install github.com/maxmind/geoipupdate/v6/cmd/geoipupdate@latest

COPY src src

RUN go build -C src -o /go/bin/app -ldflags="-w -s" -gcflags=all=-l


FROM alpine:3.18

COPY --from=builder /go/bin/geoipupdate /usr/local/bin/geoipupdate

COPY --from=builder /go/bin/app /usr/local/bin/app

COPY GeoIP.conf.sample /usr/local/etc/GeoIP.conf.sample

EXPOSE 8080

ENTRYPOINT [ "app" ]

