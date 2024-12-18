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

CMD ["bun", "run", "start:server"]