#!/usr/bin/env bash
set -e

# if flag --local exists then use go run to run site-mgr
if [ "$1" == "--local" ]; then
    COMMAND="go run ./cmd/lxd-site-mgr"
else
    COMMAND="./lxd-site-mgr"
fi

$COMMAND --state-dir ./state/dir1 external-site-join-token add site1 --expiry 24h
