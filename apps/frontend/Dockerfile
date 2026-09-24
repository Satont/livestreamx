FROM oven/bun:1.4-alpine AS base
LABEL authors="satont"
WORKDIR /app

FROM base AS build

COPY apps/frontend/package.json apps/frontend/bun.lock ./apps/frontend/
COPY apps/api/schema ./apps/api/schema/

RUN cd /app/apps/frontend && bun install --frozen-lockfile

COPY apps/frontend ./apps/frontend

RUN cd /app/apps/frontend && bun run build

FROM base
ENV NODE_ENV=production
# nitro server port, keep in sync with the traefik label (streamx-frontend)
ENV PORT=8080
# runtime image only needs the nitro output
COPY --from=build /app/apps/frontend/.output /app/.output
CMD [ "bun", "run", "/app/.output/server/index.mjs" ]
