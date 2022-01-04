#!/bin/bash

# Required envs
#   GO_REPO
#   GH_USER
#   GH_PAT
#
# Optional envs
#   GODOC_PORT - defaults to "6060"
#   REFRESH_INTERVAL_SECONDS - defaults to 180
#   GOPRIVATE_PATTERN - defaults to $(dirname $GO_REPO)/*
#   MOD_VERSION - defaults to "latest", could be versions like v1.2.3, branch names,
#     or commit hashes

GODOC_PORT=${GODOC_PORT:-"6060"}
REFRESH_INTERVAL_SECONDS=${REFRESH_INTERVAL_SECONDS:-"180"}
GOPRIVATE_PATTERN=${GOPRIVATE_PATTERN:-$(echo $(dirname $GO_REPO)'/*')}
MOD_VERSION=${MOD_VERSION:-"latest"}

export GOPRIVATE="$GOPRIVATE_PATTERN"

git config \
  --global \
  url."https://${GH_USER}:${GH_PAT}@github.com".insteadOf \
  "https://github.com"

code_path=/go/src/proj

rm -rf $code_path
mkdir -p $code_path
cd $code_path

go mod init

cat <<EOF > main.go
package main

import (
    _ "${GO_REPO}"
)

func main() {}
EOF

go get -u ${GO_REPO}@${MOD_VERSION}

godoc -index_interval ${REFRESH_INTERVAL_SECONDS}s -http :${GODOC_PORT} &
godoc_pid=$!
trap '{ kill $godoc_pid; exit 0; }' INT TERM

while /bin/true; do
	sleep $(($REFRESH_INTERVAL_SECONDS + 1)) &
	wait $!
	go get -u ${GO_REPO}@${MOD_VERSION}
done


