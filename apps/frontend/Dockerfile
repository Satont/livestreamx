FROM node:22-alpine as base
LABEL authors="satont"
WORKDIR /app

ARG PORT=3000

FROM base as build

RUN corepack enable

COPY apps/frontend/package.json ./apps/frontend/pnpm-lock.yaml /app/apps/frontend/
COPY apps/api/schema /app/apps/api/schema/

RUN cd /app/apps/frontend && pnpm install --frozen-lockfile

COPY apps/frontend /app/apps/frontend

RUN cd /app/apps/frontend && pnpm build

FROM base
ENV NODE_ENV=production
COPY --from=build /app/apps/frontend/.output /app/.output
CMD [ "node", ".output/server/index.mjs" ]