FROM golang:1.25-alpine AS builder

WORKDIR /app

# Temel araçlar
RUN apk add --no-cache git ca-certificates

# go.work ve tüm modüller için context’in kökünü kullan
COPY go.work ./
COPY ./app ./app
COPY ./packages ./packages
COPY .env /app/.env

# Modülleri indir
RUN go work sync
RUN go mod download

# Build
WORKDIR /app/app
RUN go build -o lemoras-modules .


# Use a minimal image for runtime
FROM alpine:latest

WORKDIR /app
COPY .env /app/.env
# Copy the binary from builder
COPY --from=builder /app/app/lemoras-modules .

# Copy SSL certificates if needed (optional)
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

# Set environment variables (can override with docker-compose)
ENV PORT=80
# Environment default (docker-compose override eder)
ENV DATABASE_URL=postgresql://postgres:mysecretpassword@lemoras-db:5432/zoe?sslmode=disable
ENV database_url=postgresql://postgres:mysecretpassword@lemoras-db:5432/zoe?sslmode=disable


EXPOSE 80

CMD ["./lemoras-modules"]
