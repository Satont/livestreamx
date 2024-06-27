FROM node:22-alpine as builder
LABEL authors="satont"

WORKDIR /app

RUN corepack enable

COPY apps/frontend/package.json ./apps/frontend/pnpm-lock.yaml /app/apps/frontend/
COPY apps/api/schema /app/apps/api/schema/

RUN cd /app/apps/frontend && pnpm install --frozen-lockfile

COPY apps/frontend /app/apps/frontend

ARG VITE_STREAM_URL
ENV VITE_STREAM_URL=$VITE_STREAM_URL

RUN cd /app/apps/frontend && pnpm build

FROM devforth/spa-to-http:latest
COPY --from=builder /app/apps/frontend/dist .