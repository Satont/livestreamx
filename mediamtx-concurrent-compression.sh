#!/bin/bash

ffmpeg -i rtsp://localhost:$RTSP_PORT/$MTX_PATH \
  -vf scale=w=640:h=360:force_original_aspect_ratio=decrease \
  -c:v h264 \
  -sc_threshold 0 -g 48 -keyint_min 48 \
  -profile:v main \
  -b:v 800k -maxrate 856k -bufsize 1200k -b:a 128k \
  -max_muxing_queue_size 1024 \
  -f rtsp rtsp://localhost:$RTSP_PORT/360p_$G1 &

ffmpeg -i rtsp://localhost:$RTSP_PORT/$MTX_PATH \
  -vf scale=w=842:h=480:force_original_aspect_ratio=decrease \
  -c:v h264 \
  -sc_threshold 0 -g 48 -keyint_min 48 \
  -profile:v main \
  -b:v 1400k -maxrate 1498k -bufsize 2100k -b:a 128k \
  -max_muxing_queue_size 1024 \
  -f rtsp rtsp://localhost:$RTSP_PORT/480p_$G1 &

ffmpeg -i rtsp://localhost:$RTSP_PORT/$MTX_PATH \
  -vf scale=w=1280:h=720:force_original_aspect_ratio=decrease \
  -c:v h264 \
  -sc_threshold 0 -g 48 -keyint_min 48 \
  -profile:v main \
  -b:v 2800k -maxrate 2996k -bufsize 4200k -b:a 128k \
  -max_muxing_queue_size 1024 \
  -f rtsp rtsp://localhost:$RTSP_PORT/720p_$G1 &

while true; do ffmpeg -i rtsp://localhost:$RTSP_PORT/$G1 -vframes 1 -y /thumbnails/$G1.jpg; sleep 10; done