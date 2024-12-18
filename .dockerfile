# Dockerfile
FROM oven/bun:1 as builder

WORKDIR /app
COPY . .

RUN bun install
RUN bun run build:server
RUN bun run build:web

FROM oven/bun:1-slim as runner

WORKDIR /app
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/package.json .
COPY --from=builder /app/bun.lockb .

RUN bun install --production

EXPOSE 4173 9999

# Create a startup script
RUN echo '#!/bin/sh\nbun run start:server & bun run start:web & wait' > start.sh && \
    chmod +x start.sh

CMD ["./start.sh"]