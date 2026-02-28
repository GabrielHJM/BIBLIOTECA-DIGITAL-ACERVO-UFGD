# Estágio de Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Instala dependências e cria usuário de segurança
RUN apk add --no-cache git ca-certificates && \
	addgroup -S appgroup && adduser -S appuser -G appgroup

# Gerencialmento de cache de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o código fonte
COPY . .

# Compila o binário otimizado (estático e sem símbolos de debug)
RUN CGO_ENABLED=0 GOOS=linux go build \
	-ldflags="-s -w" \
	-o main ./cmd/server/main.go

# Estágio Final
FROM alpine:latest

WORKDIR /app

# Copia certificados e usuário do builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copia apenas o binário necessário
COPY --from=builder /app/main .

# Switch para usuário não-root antes de rodar
USER appuser

EXPOSE 8080

CMD ["./main"]
