FROM golang:alpine3.24 as builder

WORKDIR /app

COPY . .

RUN go build -o log-output .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/log-output .

CMD ["./log-output"]
