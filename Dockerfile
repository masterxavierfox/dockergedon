FROM golang:1.17 AS builder

WORKDIR /app

COPY . .

RUN go build -o dockerdon

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/dockerdon .

ENTRYPOINT ["./dockerdon"]