build:
	go mod tidy
	docker build -t fetch-server --platform=linux/amd64,linux/aarch64 .

run:
	docker run -p 3000:3000 fetch-server ./fetch-server

test:
	go test ./src/tests -v