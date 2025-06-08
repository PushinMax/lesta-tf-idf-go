FROM golang:1.24.3-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o text-analyzer ./cmd/server/main.go

# --------
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/text-analyzer .

COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static


COPY configs/config.yml . 
COPY .env .
# ENV APP_PORT=8080

EXPOSE 8080

CMD ["./text-analyzer"]