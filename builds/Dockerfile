FROM golang:1.24.11-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o auth ./cmd/auth/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o blog ./cmd/blog/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o notifications ./cmd/notifications/main.go


FROM alpine:3.23 AS auth

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/auth .

EXPOSE 8080

CMD ["./auth"]


FROM alpine:3.23 AS blog

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/blog .

CMD ["./blog"]


FROM alpine:3.23 AS notifications

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/notifications .

CMD ["./notifications"]