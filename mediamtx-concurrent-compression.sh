#!/bin/bash

while true; do ffmpeg -i rtsp://localhost:$RTSP_PORT/$G1 -vframes 1 -y /thumbnails/$G1.jpg; sleep 10; done
