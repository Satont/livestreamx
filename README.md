# Development

### Requirements

* [Node.js (20+)](https://nodejs.org/en)
* [Pnpm](https://pnpm.io/)
* [Go (1.21+)](https://go.dev/)

* [Docker](https://docs.docker.com/engine/)


### Start

* Run needed services (Postgres, Mediamtx, e.t.c)
```bash
docker compose -f docker-compose.dev.yml up -d
```

* Install dependencies
```bash
cd apps/frontend
pnpm install --frozen-lockfile
```

```bash
cd apps/api
go mod download
```

* `cp .env.example .env` and fill required envs

* Run api
```bash
go run apps/api/cmd/main.go
```

* Run frontend
```bash
cd apps/frontend
pnpm dev
```

* Run stream (optional)
    * Go to OBS -> Settings -> Stream
    * Set server to `rtmp://localhost`
    * Copy stream key from `Profile` -> `Stream` from the site

### Writing migrations

* Create migration file
```bash
make create-migration name=<YOUR NAME HERE>
```

* Go to `migrations` dir and edit your newly created migration
