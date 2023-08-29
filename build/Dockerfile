FROM golang:alpine AS builder

WORKDIR /app

COPY . /app

RUN go build -o app cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/app .