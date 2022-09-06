FROM golang:1.17.11-alpine as builder

RUN apk --no-cache add tzdata

WORKDIR /build

COPY ./ ./

RUN CGO_ENABLED=0 go build -installsuffix 'static' -o /app .

FROM alpine:3.14
WORKDIR /
COPY --from=builder /app /app
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
VOLUME /config.json
EXPOSE 9015
ENTRYPOINT ["/app"]