# Stage 1: Build Frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o nexusterm main.go

# Stage 3: Final Image
FROM alpine:latest
WORKDIR /app

# Install runtime dependencies (e.g. ca-certificates for HTTPS)
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from backend builder
COPY --from=backend-builder /app/backend/nexusterm .

# Copy static files from frontend builder
# The backend expects static files in ./static relative to the binary
COPY --from=frontend-builder /app/frontend/dist ./static

# Expose port
EXPOSE 8080

# Run
CMD ["./nexusterm"]
