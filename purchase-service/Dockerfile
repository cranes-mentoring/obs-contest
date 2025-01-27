FROM golang:1.23.2 AS builder

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /purchase-service ./cmd/purchase-service/main.go


FROM alpine:3.18

RUN apk add --no-cache ca-certificates

COPY --from=builder /purchase-service /purchase-service

RUN chmod +x /purchase-service

EXPOSE 38080

CMD ["/purchase-service"]
