#Stage 1: Build Go binary
FROM golang:1.24 AS builder


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

#copy the entire project code
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /catalog-app ./internal/services/catalog-service/cmd/main.go

#Stage 2: create a small production image
FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /catalog-app .

COPY --from=builder /app/internal/services/catalog-service/config ./config

#run application
CMD ["./catalog-app"]
