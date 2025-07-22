FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy shared module first
COPY shared/ shared/

# Copy service files
COPY services/resume-processor/go.mod services/resume-processor/go.sum services/resume-processor/
WORKDIR /app/services/resume-processor
RUN go mod download

COPY services/resume-processor/ .
RUN GOWORK=off CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk update && apk --no-cache add ca-certificates wget tesseract-ocr tesseract-ocr-data-eng poppler-utils imagemagick
WORKDIR /root/

COPY --from=builder /app/services/resume-processor/main .

EXPOSE 8081
CMD ["./main"]