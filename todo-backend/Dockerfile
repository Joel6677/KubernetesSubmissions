FROM golang:alpine3.24 AS builder

WORKDIR /app

COPY . .

RUN go build -o server .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

CMD ["./server"]
