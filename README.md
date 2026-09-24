# LivestreamX

Self-hosted streaming platform: Nuxt frontend, Go API, OvenMediaEngine as the media
server (RTMP ingest, ABR transcoding, LL-HLS delivery and thumbnails).

### Requirements

* [Bun (1.4+)](https://bun.sh/)
* [Go (1.21+)](https://go.dev/)

* [Docker](https://docs.docker.com/engine/)


### Start

* Run needed services (Postgres, OvenMediaEngine, e.t.c)
```bash
docker compose -f docker-compose.dev.yml up -d
```

* Install dependencies
```bash
cd apps/frontend
bun install --frozen-lockfile
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
bun dev
```

* Run stream (optional)
    * Go to OBS -> Settings -> Stream
    * Set server to `rtmp://localhost/app`
    * Copy stream key from `Profile` -> `Stream` from the site

* Or push a test stream with ffmpeg
```bash
ffmpeg -re -f lavfi -i "testsrc2=size=1920x1080:rate=30" -f lavfi -i "sine=frequency=440" \
  -c:v libx264 -preset ultrafast -tune zerolatency -b:v 4000k -g 30 -c:a aac \
  -f flv "rtmp://localhost:1935/app/<channel>?key=<stream_key>"
```

### Streaming

Streams are ingested over RTMP and delivered as **LL-HLS** with an adaptive bitrate
ladder. OvenMediaEngine (`ome/Server.xml`) is configured with:

| Rendition | Resolution | Bitrate |
| --------- | ---------- | ------- |
| Source    | passthrough (codec untouched) | source |
| 720p      | up to 1280x720 | 2.8 Mbps |
| 480p      | up to 854x480  | 1.4 Mbps |
| 360p      | up to 640x360  | 800 kbps |

Viewers can switch quality (or leave `Auto`) in the player menu. The ladder,
encoder presets and keyframe interval are configured in `ome/Server.xml`
(`<OutputProfile>`). Transcoding runs on the CPU with the `x264` encoder and the
`faster` preset - tune it or switch to a hardware encoder for more concurrent
streams.

Streams are publicly readable; publishing is authorized by the API through
OvenMediaEngine [AdmissionWebhooks](https://docs.ovenmediaengine.com/access-control/admission-webhooks)
(`POST /streams/auth`), which validates the stream key and the channel name.

Since the source rendition is a passthrough, HDR/HEVC streams keep their original
video track and can be played by browsers with HEVC support (Safari, Chromium with
hardware decoding). The transcoded renditions are always H.264.

### Ports (development)

| Service | Address |
| ------- | ------- |
| Web / API | http://localhost:5173, http://localhost:1337 |
| RTMP ingest | rtmp://localhost:1935/app |
| LL-HLS + thumbnails | http://localhost:3334 |
| OvenMediaEngine API | http://localhost:9998 |

### Production notes

Routing is handled by Traefik through the external `traefik-public` network
(`traefik.enable=true` labels are set on `api`, `frontend` and `ome`):

| Router | Rule | Middleware | Service |
| ------ | ---- | ---------- | ------- |
| `streamx` | `Host(<domain>)` | – | frontend `:8080` |
| `streamx-http` | `Host(<domain>)` on `web` | redirect to https | frontend `:8080` |
| `streamx-api` | `Host(<domain>) && PathPrefix(/api)` | strip `/api` | api `:1337` |
| `streamx-mtx` | `Host(<domain>) && PathPrefix(/mtx)` | strip `/mtx` | ome `:3333` |

The labels assume the `web`/`websecure` entrypoints and the `le` certificate
resolver of the reference Traefik setup. The domain defaults to
`streamx.satont.dev` and can be changed with `DOMAIN` in `.env`.

The player requests `https://<domain>/mtx/app/<channel>/master.m3u8`, and
thumbnails are proxied by the API (`THUMBNAILS_URI`).

The `OME_API_ADDR` / `OME_LLHLS_ADDR` addresses are set by `docker-compose.yml`
for the api service. These envs are shared with the `ome` service and should be
changed from the defaults in the server `.env`:

```env
OME_API_ACCESS_TOKEN=change-me
OME_ADMISSION_SECRET=change-me
```

RTMP port `1935` must be exposed on the server for encoders.

### Writing migrations

* Create migration file
```bash
make create-migration name=<YOUR NAME HERE>
```

* Go to `migrations` dir and edit your newly created migration
