FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

COPY src/ .

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -o loadws-balancer .

FROM alpine:latest

RUN apk add --no-cache sqlite 

WORKDIR /app

COPY --from=builder /app/loadws-balancer .
COPY src/config.yaml ./config.yaml

VOLUME ["/app/data"]

EXPOSE 8080

CMD ["./loadws-balancer"]
