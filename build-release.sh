#!/bin/bash

set -ex

appName="chat2data"
mkdir "dist"
muslflags="--extldflags '-static -fpic'"
BASE="https://musl.nn.ci/"
FILES=(
  x86_64-linux-musl-cross
  aarch64-linux-musl-cross
  arm-linux-musleabihf-cross
#  mips-linux-musl-cross
#  mips64-linux-musl-cross
#  mips64el-linux-musl-cross
#  mipsel-linux-musl-cross
#  powerpc64le-linux-musl-cross
#  s390x-linux-musl-cross
)
for i in "${FILES[@]}"; do
  url="${BASE}${i}.tgz"
  curl -L -o "${i}.tgz" "${url}"
  sudo tar xf "${i}.tgz" --strip-components 1 -C /usr/local
done
OS_ARCHES=(
  linux-musl-amd64
  linux-musl-arm64
  linux-musl-arm
#  linux-musl-mips
#  linux-musl-mips64
#  linux-musl-mips64le
#  linux-musl-mipsle
#  linux-musl-ppc64le
#  linux-musl-s390x
)
CGO_ARGS=(
  x86_64-linux-musl-gcc
  aarch64-linux-musl-gcc
  arm-linux-musleabihf-gcc
#  mips-linux-musl-gcc
#  mips64-linux-musl-gcc
#  mips64el-linux-musl-gcc
#  mipsel-linux-musl-gcc
#  powerpc64le-linux-musl-gcc
#  s390x-linux-musl-gcc
)
for i in "${!OS_ARCHES[@]}"; do
  os_arch=${OS_ARCHES[$i]}
  cgo_cc=${CGO_ARGS[$i]}
  echo building for "${os_arch}"
  export GOOS=${os_arch%%-*}
  export GOARCH=${os_arch##*-}
  export CC=${cgo_cc}
  export CGO_ENABLED=1
  go build -o ./dist/$appName-"$os_arch" -ldflags="$muslflags" ./cmd/chat2data
done
xgo --pkg cmd/chat2data -targets=windows/amd64,darwin/amd64,darwin/arm64 -out "$appName" .
mv "$appName"-* ./dist

cd dist
mkdir compress
for i in $(find . -type f -name "$appName-linux-*"); do
  cp "$i" "$appName"
  tar -czvf compress/"$i".tar.gz "$appName"
  rm -f "$appName"
done
for i in $(find . -type f -name "$appName-darwin-*"); do
  cp "$i" "$appName"
  tar -czvf compress/"$i".tar.gz "$appName"
  rm -f "$appName"
done
for i in $(find . -type f -name "$appName-windows-*"); do
  cp "$i" "$appName".exe
  zip compress/$(echo $i | sed 's/\.[^.]*$//').zip "$appName".exe
  rm -f "$appName".exe
done