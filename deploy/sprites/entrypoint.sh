#!/usr/bin/env bash
# Runs inside the sprite. Placeholders replaced by the deploy workflow before execution:
#   __IMAGE__   → ghcr.io/<owner>/<repo>:<sha>
#   __MINUTES__ → session lifetime in minutes
set -euo pipefail

IMAGE="__IMAGE__"
MINUTES=__MINUTES__

curl -fsSL https://get.docker.com | sh

# Sprites are Firecracker microVMs without systemd — start dockerd with sudo
sudo dockerd &>/tmp/dockerd.log &

# Wait for Docker daemon socket to be ready
timeout 60 bash -c 'until docker info >/dev/null 2>&1; do sleep 2; done'

sudo docker pull "$IMAGE"
sudo docker run -d -p 8080:8080 "$IMAGE"

# Self-destruct: poweroff after session expires so compute billing stops
# (sprite hibernates — no runner kept alive in GitHub Actions)
(sleep "$((MINUTES * 60))" && sudo poweroff) &
