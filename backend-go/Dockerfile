# ---- Build stage ----
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server

# ---- Run stage ----
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/server .
COPY backend-go.env .

EXPOSE 8080

CMD ["./server"]
