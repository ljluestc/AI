# Build stage for frontend
FROM node:16-alpine AS frontend-build
WORKDIR /app/frontend
COPY web/package*.json ./
RUN npm install
COPY web/ ./
RUN npm run build

# Build stage for backend
FROM golang:1.24-alpine AS backend-build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o teathis-server

# Final stage
FROM alpine:3.18
RUN apk --no-cache add ca-certificates git

WORKDIR /app

# Copy the compiled binaries and assets
COPY --from=backend-build /app/teathis-server .
COPY --from=frontend-build /app/frontend/build ./web/build

# Create workspace directory
RUN mkdir -p /app/workspace

# Expose the application port
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release
ENV PORT=8080

# Run the application
CMD ["./teathis-server"]
