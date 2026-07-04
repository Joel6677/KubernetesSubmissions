FROM golang:alpine3.24 as builder

WORKDIR /app

COPY . .

RUN go build -o ping-pong .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ping-pong .

CMD ["./ping-pong"]
