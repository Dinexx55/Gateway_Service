FROM golang:1.21.0 AS builder

WORKDIR /gateway/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main cmd/app/main.go

FROM golang:1.21.0

COPY --from=builder /gateway/app/main /
COPY --from=builder /gateway/app/configs/ /configs

EXPOSE 8080

CMD /main