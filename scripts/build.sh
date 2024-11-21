#!/bin/bash
docker build -t gobuilder .
mkdir -p dist
docker run --rm -v "$(pwd)/dist:/dist" gobuilder
echo "Build complete! Binary is in the ./dist directory"