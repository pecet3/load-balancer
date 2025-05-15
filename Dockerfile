FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY src/ .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o loadws-balancer .

FROM scratch

COPY --from=builder /app/loadws-balancer /loadws-balancer
COPY src/cfg/config.yaml /cfg/config.yaml

VOLUME ["/cfg"]

EXPOSE 8080

CMD ["/loadws-balancer"]
