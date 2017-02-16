test:
	go test -v -cover -race ./...

lint:
	golint ./...
	go vet ./...
	errcheck ./...

run:
	go run go-carpet.go ast.go utils.go terminal_unix.go -256colors

update-from-github:
	go get -u github.com/msoap/go-carpet

gometalinter:
	gometalinter --vendor --cyclo-over=20 --line-length=150 --dupl-threshold=150 --min-occurrences=2 --enable=misspell --deadline=10m ./...

generate-manpage:
	docker run -it --rm -v $$PWD:/app -w /app ruby-ronn sh -c 'cat README.md | grep -v "^\[" | grep -v Screenshot > go-carpet.md; ronn go-carpet.md; mv ./go-carpet ./go-carpet.1; rm ./go-carpet.html ./go-carpet.md'

create-debian-amd64-package:
	GOOS=linux GOARCH=amd64 go build -ldflags="-w" -o go-carpet
	set -e ;\
	TAG_NAME=$$(git tag 2>/dev/null | grep -E '^[0-9]+' | tail -1) ;\
	docker run -it --rm -v $$PWD:/app -w /app -e TAG_NAME=$$TAG_NAME ruby-fpm sh -c 'fpm -s dir -t deb --name go-carpet -v $$TAG_NAME ./go-carpet=/usr/bin/ ./go-carpet.1=/usr/share/man/man1/ LICENSE=/usr/share/doc/go-carpet/copyright README.md=/usr/share/doc/go-carpet/'
	rm go-carpet
