# Build stage: compile the Vue SPA.
FROM node:20-alpine AS builder
WORKDIR /app
COPY vue-app/package.json vue-app/package-lock.json ./
RUN npm ci
COPY vue-app/ ./
ENV VITE_API_BASE_URL=/api
RUN npm run build

# Runtime stage: nginx serves the SPA. Backend routing is handled by Traefik;
# nginx only does SPA fallback so client-side routes don't 404 on refresh.
FROM nginx:1.27-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
