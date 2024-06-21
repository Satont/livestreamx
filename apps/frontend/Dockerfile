FROM node:22-alpine as builder
LABEL authors="satont"

WORKDIR /app

RUN corepack enable

COPY apps/frontend/package.json ./apps/frontend/pnpm-lock.yaml /app/apps/frontend/
COPY apps/api/schema /app/apps/api/schema/

RUN cd /app/apps/frontend && pnpm install

COPY apps/frontend ./

RUN cd /app/apps/frontend && pnpm build

FROM devforth/spa-to-http:latest
COPY --from=builder /app/dist /app