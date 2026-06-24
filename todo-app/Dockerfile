FROM golang:alpine3.24 as builder

WORKDIR /app

COPY . .

RUN go build -o todo-app .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/todo-app .

CMD ["./todo-app"]
