FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app .
EXPOSE 8080
CMD ["./main"]