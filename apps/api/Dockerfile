FROM golang:1-alpine3.20 as builder
LABEL authors="satont"

WORKDIR /app

COPY go.work go.work.sum ./
COPY apps/api/go.mod apps/api/go.sum /app/apps/api/

RUN go mod download

COPY apps/api ./apps/api

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /tmp/api ./apps/api/cmd/main.go

FROM alpine:3.20
WORKDIR /app
COPY migrations /app/migrations
COPY --from=builder /tmp/api /app/api
CMD ["/app/api"]