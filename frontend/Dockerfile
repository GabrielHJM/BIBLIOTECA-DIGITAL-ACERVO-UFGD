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

# Configuração customizada do Nginx para suporte a SPA e Proxy
COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY --from=build-stage /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
