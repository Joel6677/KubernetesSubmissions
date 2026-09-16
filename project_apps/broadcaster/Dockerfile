FROM golang:alpine3.24 AS builder
WORKDIR /app
COPY . .
RUN go build -o broadcaster .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/broadcaster .
CMD ["./broadcaster"]
