FROM golang:1.23.2 AS builder

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /auth-service ./cmd/auth-service/main.go


FROM alpine:3.18

RUN apk add --no-cache ca-certificates

COPY --from=builder /auth-service /auth-service

RUN chmod +x /auth-service

EXPOSE 18080

CMD ["/auth-service"]
