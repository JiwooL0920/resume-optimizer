FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOWORK=off CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk update && apk --no-cache add ca-certificates wget tesseract-ocr tesseract-ocr-data-eng poppler-utils imagemagick
WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8081
CMD ["./main"]