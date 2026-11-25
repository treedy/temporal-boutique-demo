#!/usr/bin/env bash
# Build and push the checkout worker or service container.
set -euo pipefail

REGISTRY="us-east4-docker.pkg.dev/treedy-temporal-store-demo/containers"

case "${1:-}" in
  worker)
    NAME="checkoutworker"
    DOCKERFILE="temporal.worker.Dockerfile"
    ;;
  service)
    NAME="checkoutservice"
    DOCKERFILE="Dockerfile"
    ;;
  *)
    echo "Usage: $0 worker|service" >&2
    exit 1
    ;;
esac

cd "$(dirname "$0")"

TAG="$(date +%Y%m%d-%H%M%S)"
IMAGE="${REGISTRY}/${NAME}:${TAG}"
REPLACE_IMG="${REGISTRY}/${NAME}:v0.10.7-tmprl"

docker build -f "$DOCKERFILE" -t "$IMAGE" .
docker tag "$IMAGE" "$REPLACE_IMG"
docker push "$REPLACE_IMG"

#docker push "$IMAGE"
# Also tag with the replacement tag

echo "Pushed $REPLACE_IMG"
