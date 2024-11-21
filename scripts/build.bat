@echo off
docker build -t gobuilder .
if not exist dist mkdir dist
docker run --rm -v %cd%/dist:/dist gobuilder
echo Build complete! Binary is in the ./dist directory