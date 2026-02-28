# Estágio de Build
FROM node:18-alpine AS build-stage

WORKDIR /app

# Cache de dependências
COPY package*.json ./
RUN npm ci --quiet

COPY . .
RUN npm run build

# Estágio de Produção
FROM nginx:stable-alpine AS production-stage

# Configuração customizada do Nginx para suporte a SPA (fallback para index.html)
RUN echo 'server { \
	listen 80; \
	location / { \
	root /usr/share/nginx/html; \
	index index.html; \
	try_files $uri $uri/ /index.html; \
	} \
	}' > /etc/nginx/conf.d/default.conf

COPY --from=build-stage /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
