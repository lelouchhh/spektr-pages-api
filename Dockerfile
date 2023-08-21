
FROM golang:1.19.0

WORKDIR /spektr-pages-api
RUN go install github.com/cosmtrek/air@latest

COPY . .