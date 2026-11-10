FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod ./
COPY main.go ./

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/test-server .

FROM alpine:3.22

RUN adduser -D -H -u 10001 appuser

COPY --from=builder /out/test-server /usr/local/bin/test-server

USER appuser
EXPOSE 8081

ENTRYPOINT ["/usr/local/bin/test-server"]
