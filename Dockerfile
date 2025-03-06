FROM docker.io/golang:1.24 AS builder
WORKDIR /app

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o oteller main.go 

# Final stage
FROM docker.io/alpine:latest
RUN mkdir /app && adduser -h /app -D oteller
WORKDIR /app
COPY --chown=oteller:oteller --from=builder /app/oteller .

# Expose the application port
EXPOSE 8080

# Run the application
USER oteller
ENTRYPOINT ["/app/oteller"] 