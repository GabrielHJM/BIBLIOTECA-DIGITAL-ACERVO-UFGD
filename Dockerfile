# Stage 1: Build Frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.22-alpine AS backend-builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
# Copy the built frontend to a "dist" folder inside the backend build context
# so the Go server can serve it.
COPY --from=frontend-builder /app/frontend/dist ./dist
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# Stage 3: Final Image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Copy binary and dist folder
COPY --from=backend-builder /app/backend/main .
COPY --from=backend-builder /app/backend/dist ./dist

EXPOSE 8080
CMD ["./main"]
