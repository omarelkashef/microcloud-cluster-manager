#!/bin/sh

# The script requires:
# - rockcraft
# - yq
# - docker

set -e

rockcraft pack -v

VERSION=$(awk -F': ' '/^version:/ {print $2}' rockcraft.yaml | tr -d '"')
ROCK="lxd-cluster-manager_${VERSION}_amd64.rock"

# $IMAGE is env var set by skaffold for custom artifact builders
rockcraft.skopeo --insecure-policy copy \
    oci-archive:$ROCK \
    docker-daemon:$IMAGE

# remove rock after copying over to docker daemon
rm $ROCK