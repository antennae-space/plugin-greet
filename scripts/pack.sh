#!/bin/sh
# Builds the release archives of the plugin, one per architecture: dist/greet-<arch>.tar.gz
# holds plugin.json, web/, i18n/ and bin/greet (the backend, static). The version of the tag
# (v1.2.3) must be the one of plugin.json: the IDE checks it too.
set -eu
cd "$(dirname "$0")/.."
version=$(sed -n 's/^  "version": "\(.*\)",$/\1/p' plugin.json)
if [ -n "${TAG:-}" ] && [ "$TAG" != "v$version" ]; then
  echo "the tag $TAG is not the version $version of plugin.json" >&2
  exit 1
fi
rm -rf dist
for arch in amd64 arm64; do
  dir=dist/linux-$arch
  mkdir -p "$dir/bin"
  CGO_ENABLED=0 GOOS=linux GOARCH=$arch go build -trimpath -ldflags "-s -w" -o "$dir/bin/greet" ./backend
  cp -r plugin.json web i18n "$dir/"
  tar -czf "dist/greet-linux-$arch.tar.gz" -C "$dir" .
done
ls -l dist/*.tar.gz
