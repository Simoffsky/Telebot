.PHONY: clean build-darwin build-windows build-linux
clean:
	rm -rf ./bin/*.out
	rm -rf ./bin/*.exe

run: 
	go run ./cmd/vladOS/main.go
build: build-darwin build-windows build-linux

run-weather:
	go run ./cmd/weatherparser/main.go

build-darwin:
	GOOS=darwin go build  -o ./bin/darwin.out -v ./cmd/vladOS/main.go
build-windows:
	GOOS=windows go build -o ./bin/windows.exe -v ./cmd/vladOS/main.go
build-linux:
	GOOS=linux go build -o ./bin/linux.out -v ./cmd/vladOS/main.go
