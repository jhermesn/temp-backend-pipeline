#!/usr/bin/env bash
# EC2 user-data. Placeholders replaced by the deploy workflow before launch:
#   __MINUTES__ → lifetime in minutes
#   __IMAGE__   → ghcr.io/<owner>/<repo>:<sha>
set -euo pipefail

MINUTES=__MINUTES__
IMAGE="__IMAGE__"

# Self-destruct: belt-and-suspenders alongside the Actions job teardown
shutdown -P "+$MINUTES" &

# Install Docker (Amazon Linux 2023)
dnf install -y docker
systemctl start docker
systemctl enable docker

# Wait for Docker daemon
timeout 30 bash -c 'until docker info >/dev/null 2>&1; do sleep 2; done'

docker pull "$IMAGE"
docker run -d -p 80:8080 "$IMAGE"
