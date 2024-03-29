FROM golang:1.20-alpine as builder

WORKDIR /app
COPY ./ /app
RUN apk add --no-cache make gcc musl-dev linux-headers
RUN go mod tidy
RUN go build -o ./builds/linux/explorer ./cmd/explorer.go

FROM alpine:3.17

COPY --from=builder /app/builds/linux/explorer /usr/bin/explorer
RUN addgroup minteruser && adduser -D -h /minter -G minteruser minteruser
USER minteruser
WORKDIR /minter
ENTRYPOINT ["/usr/bin/explorer"]
CMD ["explorer"]