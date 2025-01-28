# Fetch Rewards App
## Requirements
- `golang` version 1.23.3
- Docker must be installed

## Usage
1. Build the docker image via `make build`
2. Start the server via `make run`
3. Optionally run tests via `make test` once the server is running

## Notes
- Docker image supports `linux/amd64` and `linux/aarch64` platforms. Feel free to update the Dockerfile to support more platforms