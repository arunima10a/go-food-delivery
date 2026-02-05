#Stage 1: Build Go binary
FROM golang:1.24 AS builder


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

#copy the entire project code
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /search-app ./internal/services/search-service/cmd/main.go

#Stage 2: create a small production image
FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /search-app .

COPY --from=builder /app/internal/services/search-service/config ./config

#run application
CMD ["./search-app"]
