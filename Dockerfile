# ---- Builder stage ----
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ---- Final stage ----
FROM alpine:latest

RUN apk add --no-cache tzdata ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080
# EXPOSE 50051

CMD ["./main"]