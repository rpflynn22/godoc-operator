#!/bin/bash

# Required envs
#   GO_REPO
#   GH_PAT
#
# Optional envs
#   GODOC_PORT - defaults to "6060"
#   GIT_DEFAULT_BRANCH - defaults to "main"
#   REFRESH_INTERVAL_SECONDS - defaults to 180

GODOC_PORT=${GODOC_PORT:-"6060"}
GIT_DEFAULT_BRANCH=${GIT_DEFAULT_BRANCH:-"main"}
REFRESH_INTERVAL_SECONDS=${REFRESH_INTERVAL_SECONDS:-"180"}

full_path=/go/src/$GO_REPO

rm -rf $full_path
git clone https://$GH_PAT@$GO_REPO $full_path >/dev/null 2>&1

godoc -index_interval ${REFRESH_INTERVAL_SECONDS}s -http :${GODOC_PORT} -goroot /go &
godoc_pid=$!
trap '{ kill $godoc_pid; exit 0; }' INT TERM

cd $full_path

while /bin/true; do
	sleep $(($REFRESH_INTERVAL_SECONDS + 1)) &
	wait $!
	git fetch origin >/dev/null 2>&1
	git reset --hard origin/${GIT_DEFAULT_BRANCH} >/dev/null 2>&1
done


