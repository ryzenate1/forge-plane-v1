FROM node:20-alpine AS deps
WORKDIR /src
COPY package.json package-lock.json ./
RUN npm ci && npm cache clean --force

FROM node:20-alpine AS build
WORKDIR /src
ARG API_INTERNAL_URL=http://api:8080
ENV API_INTERNAL_URL=$API_INTERNAL_URL
ARG NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
ARG NEXT_PUBLIC_DEMO_MODE=false
ENV NEXT_PUBLIC_DEMO_MODE=$NEXT_PUBLIC_DEMO_MODE
ARG NEXT_PUBLIC_SFTP_HOST=127.0.0.1
ENV NEXT_PUBLIC_SFTP_HOST=$NEXT_PUBLIC_SFTP_HOST
ARG NEXT_PUBLIC_SFTP_PORT=2222
ENV NEXT_PUBLIC_SFTP_PORT=$NEXT_PUBLIC_SFTP_PORT
ARG NEXT_PUBLIC_SFTP_USER=game
ENV NEXT_PUBLIC_SFTP_USER=$NEXT_PUBLIC_SFTP_USER
COPY --from=deps /src/node_modules ./node_modules
COPY . .
RUN npm run build

FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
ENV HOSTNAME=0.0.0.0
ENV PORT=3000
RUN apk add --no-cache curl
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=build /src/.next/standalone ./
COPY --from=build /src/.next/static ./.next/static
RUN chown -R appuser:appgroup /app
USER appuser
EXPOSE 3000
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 CMD ["curl", "-f", "http://localhost:3000"]
CMD ["node", "server.js"]
