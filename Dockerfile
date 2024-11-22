FROM golang:1.16.3-alpine3.13 AS moduler
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

FROM golang:1.16.3-alpine3.13 AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linus GOARCH=arm64 \
    go build -o /bin/app ./cmd/app

FROM scratch
COPY --from=builder /app/config /config
COPY --from=builder /app/migrations /migrations
COPY --from=builder /bin/app /app
CMD ["/app"]