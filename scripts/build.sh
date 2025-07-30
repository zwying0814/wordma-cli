#!/bin/bash

# Build script for local development
# Usage: ./scripts/build.sh [version]

set -e

VERSION=${1:-"dev"}
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS="-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

echo "Building wordma CLI..."
echo "Version: ${VERSION}"
echo "Build Time: ${BUILD_TIME}"
echo "Git Commit: ${GIT_COMMIT}"
echo ""

# Create dist directory
mkdir -p dist

# Build for multiple platforms
platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$platform"
    
    output_name="wordma-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    echo "Building for ${GOOS}/${GOARCH}..."
    
    env GOOS="$GOOS" GOARCH="$GOARCH" go build \
        -ldflags="$LDFLAGS" \
        -o "dist/${output_name}" \
        .
    
    if [ $? -ne 0 ]; then
        echo "Failed to build for ${GOOS}/${GOARCH}"
        exit 1
    fi
done

echo ""
echo "Build completed successfully!"
echo "Binaries are available in the dist/ directory:"
ls -la dist/