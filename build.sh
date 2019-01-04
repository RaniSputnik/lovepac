#!/bin/bash

BUILD_DIR="build"
if [ ! -d "$BUILD_DIR" ]; then
    mkdir $BUILD_DIR
fi

build()
{
    TARGET="build/lovepac-$GOOS-$GOARCH"
    go build -o "${TARGET}"
}

jobs=(
    "GOOS=windows GOARCH=amd64  go build -o build/lovepac-windows-amd64.exe" \
    "GOOS=darwin GOARCH=amd64   build" \
    "GOOS=linux GOARCH=amd64    build" \
    "GOOS=linux GOARCH=arm      build" \
    "GOOS=linux GOARCH=ppc64le  build" \
    "GOOS=linux GOARCH=s390x    build" \
)

echo "Building binaries for all platforms"
for j in "${jobs[@]}" 
do 
    eval $j &
done
wait
echo "Build completed"