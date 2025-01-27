FROM golang:1.23.2 AS builder

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /purchase-processor ./cmd/purchase-processor/main.go


FROM alpine:3.18

RUN apk add --no-cache ca-certificates

COPY --from=builder /purchase-processor /purchase-processor

RUN chmod +x /purchase-processor

EXPOSE 48080

CMD ["/purchase-processor"]
