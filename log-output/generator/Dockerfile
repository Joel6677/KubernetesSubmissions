FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o generator .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/generator .
CMD ["./generator"]
