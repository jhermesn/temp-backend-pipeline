#!/usr/bin/env bash
# Runs inside the sprite. Placeholder replaced by the deploy workflow before execution:
#   __IMAGE__ → ghcr.io/<owner>/<repo>:<sha>
set -euo pipefail

IMAGE="__IMAGE__"

curl -fsSL https://get.docker.com | sh
dockerd &>/tmp/dockerd.log &

# Wait for Docker daemon socket to be ready
timeout 30 bash -c 'until docker info >/dev/null 2>&1; do sleep 2; done'

docker pull "$IMAGE"
docker run -d -p 80:8080 "$IMAGE"