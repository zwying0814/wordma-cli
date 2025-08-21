#!/bin/bash

# Build script for local development
# Usage: ./scripts/build.sh [version]

set -e

VERSION=${1}
if [ -z "$VERSION" ]; then
    # Try to get version from git tag
    VERSION=$(git describe --tags --exact-match HEAD 2>/dev/null || echo "")
    if [ -z "$VERSION" ]; then
        GIT_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
        if [ -z "$GIT_TAG" ]; then
            VERSION="dev"
        else
            VERSION="${GIT_TAG}-dev"
        fi
    fi
fi
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
    
    # 设置临时文件名和最终文件名
    temp_name="temp-wordma-${GOOS}-${GOARCH}"
    final_name="wordma"
    archive_name="wordma-${GOOS}-${GOARCH}"
    
    if [ "$GOOS" = "windows" ]; then
        temp_name="${temp_name}.exe"
        final_name="${final_name}.exe"
        archive_name="${archive_name}.zip"
    else
        archive_name="${archive_name}.tar.gz"
    fi
    
    echo "Building for ${GOOS}/${GOARCH}..."
    
    env GOOS="$GOOS" GOARCH="$GOARCH" go build \
        -ldflags="$LDFLAGS" \
        -o "dist/${temp_name}" \
        .
    
    if [ $? -ne 0 ]; then
        echo "Failed to build for ${GOOS}/${GOARCH}"
        exit 1
    fi
    
    # 创建压缩包
    cd dist
    if [ "$GOOS" = "windows" ]; then
        # 重命名并创建zip
        mv "$temp_name" "$final_name"
        zip "$archive_name" "$final_name"
        rm "$final_name"
    else
        # 重命名并创建tar.gz
        mv "$temp_name" "$final_name"
        tar -czf "$archive_name" "$final_name"
        rm "$final_name"
    fi
    cd ..
    
    echo "Created $archive_name"
done

echo ""
echo "Build completed successfully!"
echo "Archives are available in the dist/ directory:"
ls -la dist/*.zip dist/*.tar.gz 2>/dev/null || true