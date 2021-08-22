FROM golang:1.16 AS builder

WORKDIR /go/src/github.com/tomoconnor/apcupsd_exporter
COPY . /go/src/github.com/tomoconnor/apcupsd_exporter/
RUN go build /go/src/github.com/tomoconnor/apcupsd_exporter/cmd/apcupsd_exporter



FROM alpine:latest
RUN apk add libc6-compat
WORKDIR /opt/apcupsd_exporter
COPY --from=builder /go/src/github.com/tomoconnor/apcupsd_exporter/apcupsd_exporter /opt/apcupsd_exporter/apcupsd_exporter
EXPOSE 9162
CMD ["/opt/apcupsd_exporter/apcupsd_exporter"]