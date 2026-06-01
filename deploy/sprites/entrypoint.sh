#!/usr/bin/env bash
# Runs inside the sprite. Expects IMAGE env var to be set by the caller.
set -euo pipefail

: "${IMAGE:?IMAGE env var must be set}"

curl -fsSL https://get.docker.com | sh
systemctl start docker

# Wait for Docker daemon
timeout 30 bash -c 'until docker info >/dev/null 2>&1; do sleep 2; done'

docker pull "$IMAGE"
docker run -d -p 80:8080 "$IMAGE"
