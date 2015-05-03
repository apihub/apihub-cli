help:
	@echo '    save-deps ................ generates the Godeps folder'
	@echo '    test ..................... runs tests'

save-deps:
	$(GOPATH)/bin/godep save ./...

test:
	go test ./...