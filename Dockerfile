# ---- build stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o service ./cmd

# ---- runtime stage ----
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/docs ./docs
COPY --from=builder /app/service .
COPY --from=builder /app/assets ./assets

EXPOSE 8080

CMD ["./service"]