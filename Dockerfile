FROM docker.io/golang:1.24 AS builder
WORKDIR /app

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o rolldice main.go

# Final stage
FROM docker.io/alpine:latest
RUN mkdir /app && adduser -h /app -D rolldice
WORKDIR /app
COPY --chown=rolldice:rolldice --from=builder /app/rolldice .

# Expose the application port
EXPOSE 8080

# Run the application
USER rolldice
ENTRYPOINT ["/app/rolldice"] 