build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/s3rw s3readwrite/main.go
