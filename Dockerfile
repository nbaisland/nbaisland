FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o nbaisland-backend ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/nbaisland-backend .

COPY --from=builder /app/db/migrations ./db/migrations

RUN mkdir -p /app/logs

EXPOSE 8080

CMD ["./nbaisland-backend"]
