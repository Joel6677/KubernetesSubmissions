FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o greeter .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/greeter .
CMD ["./greeter"]

