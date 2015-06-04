include go.mk

.PHONY: help
help:
	$(info | build ........... builds the binary                                  |)
	$(info | clean ........... remove all untracked and temporary files           |)
	$(info | test ............ run all tests                                      |)
	$(info | race ............ run all tests with race detection enabled          |)
	$(info | cover ........... run all tests with coverage enabled                |)
	$(info | vet ............. run the go vet tool                                |)
	$(info | lint ............ run the golint tool                                |)
	$(info | save-deps ....... generates the Godeps folder                        |)

build: gomkbuild

# Kept for backwards compatibility. User can call `make savegodeps` directly
.PHONY: save-deps
save-deps: savegodeps

.PHONY: clean
clean: gomkclean
