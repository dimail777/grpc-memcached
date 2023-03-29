.PHONY: app-compile app-build-bin app-build-docker

app-compile:
	echo "Compiling proto files"
    protoc --go_out=internal --go_opt=paths=import --go-grpc_out=internal --go-grpc_opt=paths=import internal/proto/app.proto

app-build-bin:
	echo "Building binaries for all required OS and Platforms"
	GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -o cmd/server/bin/app-linux-arm cmd/server/main.go
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o cmd/server/bin/app-linux-386 cmd/server/main.go

app-build-docker:
	echo "Building docker image"
	docker build -t grpc.com/memcached .

e2e-test-bin:
	echo "Building e2e binaries for all required OS and Platforms"
	GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -o cmd/e2e/bin/e2e-linux-arm cmd/server/main.go
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o cmd/e2e/bin/e2e-linux-386 cmd/server/main.go

app-build: app-compile app-build-bin app-build-docker
