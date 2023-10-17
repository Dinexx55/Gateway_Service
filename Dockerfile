FROM golang:1.21.0 AS builder

WORKDIR /gateway/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o gateway-service cmd/app/main.go

FROM debian:bookworm AS runner

WORKDIR /usr/bin

COPY --from=builder /gateway/app/gateway-service .
COPY --from=builder /gateway/app/configs/ /usr/bin/configs

EXPOSE 8081

ENTRYPOINT ["gateway-service"]