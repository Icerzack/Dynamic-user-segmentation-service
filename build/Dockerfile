FROM golang:alpine AS builder

WORKDIR /app

COPY . /app
COPY /docs /app/docs

RUN go build -o avito-app cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/avito-app .
COPY --from=builder /app/docs /app/docs