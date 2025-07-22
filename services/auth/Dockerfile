FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy shared module first
COPY shared/ shared/

# Copy service files
COPY services/auth/go.mod services/auth/go.sum services/auth/
WORKDIR /app/services/auth
RUN go mod download

COPY services/auth/ .
RUN GOWORK=off CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates wget
WORKDIR /root/

COPY --from=builder /app/services/auth/main .

EXPOSE 8080
CMD ["./main"]