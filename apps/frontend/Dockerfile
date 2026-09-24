FROM node:22-alpine as base
LABEL authors="satont"
WORKDIR /app

FROM base as build

# pin pnpm: the lockfile is v9, modern pnpm releases skip dependency
# build scripts (esbuild, sharp) and fail the install
RUN npm install -g pnpm@9.15.9

COPY apps/frontend/package.json ./apps/frontend/pnpm-lock.yaml /app/apps/frontend/
COPY apps/api/schema /app/apps/api/schema/

RUN cd /app/apps/frontend && pnpm install --frozen-lockfile

COPY apps/frontend /app/apps/frontend

RUN cd /app/apps/frontend && pnpm build

FROM base
ENV NODE_ENV=production
# nitro server port, keep in sync with the traefik label (streamx-frontend)
ENV PORT=8080
COPY --from=build /app/apps/frontend/.output /app/.output
CMD [ "node", ".output/server/index.mjs" ]