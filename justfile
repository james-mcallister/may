# just is a task runner used to save and run project specific commands.

set shell := ["bash", "-uc"]

VERSION := "0.0.1"
export CGO_ENABLED := "1"
export CC := "zig cc -target x86_64-linux-musl"
export GOOS := "linux"
export GOARCH := "amd64"

# list available recipes
default:
  @just --list

# remove previously generated assets
clean:
  rm -rf frontend/dist
  rm -rf bin

# create a github release
release: tag
  gh release create {{VERSION}} --latest --notes-file CHANGELOG.md --target main --title "v{{VERSION}}" --verify-tag ./dist/*

# create and push a git tag for the currently configured version (in justfile)
tag:
  git tag -a -s {{VERSION}} -m "release new version"
  git push origin --tags

# serve frontend index.html file
serve:
  python3 -m http.server -b 127.0.0.1 -d frontend/dist 12345

# build frontend using eslint updating assets directory
build-frontend: clean
  mkdir -p frontend/dist
  cp frontend/index.html frontend/favicon.ico frontend/dist/.
  ./frontend/node_modules/.bin/esbuild frontend/src/main.js --bundle --minify --outfile=./frontend/dist/assets/bundle.js

# build frontend using eslint updating assets directory
build-frontend-dev: clean
  mkdir -p frontend/dist
  cp frontend/index.html frontend/favicon.ico frontend/dist/.
  ./frontend/node_modules/.bin/esbuild frontend/src/main.js --bundle --outfile=./frontend/dist/assets/bundle.js

# build the backend. Includes the frontend assets built into a production binary
build: build-frontend-dev
  go build -tags "linux fts5 foreign_keys json" -ldflags "-s -w" -o ./bin/may main.go domain.go
