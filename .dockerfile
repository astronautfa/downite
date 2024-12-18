FROM oven/bun:1 as builder

WORKDIR /app
COPY . .

# Install dependencies
RUN bun install

# Create a .env file with the base URL configuration
RUN echo "PUBLIC_URL=${PUBLIC_URL:-/}" > .env

# Build the web application
RUN bun run build:web

FROM oven/bun:1-slim as runner

WORKDIR /app
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/package.json .
COPY --from=builder /app/bun.lockb .

RUN bun install --production

EXPOSE 4173

CMD ["bun", "run", "start:web"]