test:
	go test -cover ./...

lint:
	golint ./...
	go vet ./...

run:
	go run go-carpet.go terminal_unix.go -256colors

update-from-github:
	go get -u github.com/msoap/go-carpet
