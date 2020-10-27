FROM golang:1.12-alpine

RUN apk add --no-cache git make

ARG transactions_mongodb_url
ENV TRANSACTIONS_MONGODB_URL=$transactions_mongodb_url

ARG transactions_mongodb_database
ENV TRANSACTIONS_MONGODB_DATABASE=$transactions_mongodb_database

ARG log_level
ENV LOG_LEVEL=$log_level

ARG sender_email
ENV SENDER_EMAIL=$sender_email

ARG receiver_email
ENV RECEIVER_EMAIL=$receiver_email

WORKDIR /app
COPY . .
RUN GO111MODULE=on go build

ENTRYPOINT ["/app/accounts-statistics-tool"]
