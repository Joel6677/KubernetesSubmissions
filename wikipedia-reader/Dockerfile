FROM alpine:latest

RUN apk add --no-cache curl bash

WORKDIR /app

COPY create-todo.sh .

RUN chmod +x ./create-todo.sh

CMD ["./create-todo.sh"]
