build:
	@mkdir -p bin
	@go build -o bin/fs .

run: build
	@bin/fs

clean:
	@rm -f bin/fs

test:
	@go test ./... -v