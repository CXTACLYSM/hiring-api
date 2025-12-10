FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o api ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o worker ./cmd/worker/main.go

FROM alpine:3.20 AS api

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]

FROM alpine:3.20 AS worker

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/worker .

CMD ["./worker"]