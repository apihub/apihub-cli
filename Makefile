help:
	@echo '    build .................... builds the binary'
	@echo '    save-deps ................ generates the Godeps folder'
	@echo '    test ..................... runs tests'

build:
	go build .

save-deps:
	$(GOPATH)/bin/godep save ./...

test:
	go test ./...
