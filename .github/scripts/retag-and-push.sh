#!/usr/bin/env bash
set -e

FROM="$1"
TO="$2"

echo "Pushing ${TO}"
docker tag "${FROM}" "${TO}"
docker push "${TO}"
