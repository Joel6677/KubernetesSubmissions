FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o reader .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/reader .
CMD ["./reader"]
