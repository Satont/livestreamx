FROM bluenviron/mediamtx:1.8.3-ffmpeg

RUN apk add --no-cache bash

RUN mkdir -p /opt/mediamtx/scripts
COPY mediamtx-concurrent-compression.sh /opt/mediamtx/scripts
RUN chmod +x /opt/mediamtx/scripts/mediamtx-concurrent-compression.sh