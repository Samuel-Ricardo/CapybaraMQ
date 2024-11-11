
FROM golang:1.22.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go test ./test/unit/ && \
  go test ./test/integration/ && \
  go test ./test/e2e/ && \
  go build -o main.exe ./cmd


FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/.env .
COPY --from=builder /app/main.exe .

EXPOSE 8080

CMD ["./main.exe"]

