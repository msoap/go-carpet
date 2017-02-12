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
