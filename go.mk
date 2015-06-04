# go.mk
#
# Copyright (c) 2015, Herbert G. Fischer
# All rights reserved.
#
# Redistribution and use in source and binary forms, with or without
# modification, are permitted provided that the following conditions are met:
#     * Redistributions of source code must retain the above copyright
#       notice, this list of conditions and the following disclaimer.
#     * Redistributions in binary form must reproduce the above copyright
#       notice, this list of conditions and the following disclaimer in the
#       documentation and/or other materials provided with the distribution.
#     * Neither the name of the organization nor the
#       names of its contributors may be used to endorse or promote products
#       derived from this software without specific prior written permission.
#
# THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
# ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
# WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
# DISCLAIMED. IN NO EVENT SHALL HERBERT G. FISCHER BE LIABLE FOR ANY
# DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
# (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
# LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
# ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
# (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
# SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

APPBIN      := $(shell basename $(PWD))
GOSOURCES   := $(shell find . -type f -name '*.go' ! -path '*Godeps/_workspace*')
GOPKGS      := $(shell go list ./...)
GOPKG       := $(shell go list)
COVERAGEOUT := coverage.out
COVERAGETMP := coverage.tmp
GODEPPATH   := $(PWD)/Godeps/_workspace
LOCALGOPATH := $(GODEPPATH):$(GOPATH)
ORIGGOPATH  := $(GOPATH)
GOMKVERSION := 0.5.0

ifndef GOBIN
export GOBIN := $(GOPATH)/bin
endif

ifndef GOPATH
$(error ERROR!! GOPATH must be declared. Check [http://golang.org/doc/code.html#GOPATH])
else
export GOPATH=$(LOCALGOPATH)
endif

ifeq ($(shell go list ./... | grep -q '^_'; echo $$?), 0)
$(error ERROR!! This directory should be at $(GOPATH)/src/$(REPO)]
endif

##########################################################################################
## Project targets
##########################################################################################

$(APPBIN): gomkbuild

##########################################################################################
## Main targets
##########################################################################################

.PHONY: gomkbuild
gomkbuild:  $(GOSOURCES) ; @go build

.PHONY: gomkxbuild
gomkxbuild: ; $(GOX)

.PHONY: gomkclean
gomkclean:
	@rm -vf $(APPBIN)_*_386 $(APPBIN)_*_amd64 $(APPBIN)_*_arm $(APPBIN)_*.exe
	@rm -vf $(COVERAGEOUT) $(COVERAGETMP)
	@go clean

##########################################################################################
## Go tools
##########################################################################################

GOTOOLDIR := $(shell go env GOTOOLDIR)
BENCHCMP  := $(GOTOOLDIR)/benchcmp
CALLGRAPH := $(GOTOOLDIR)/callgraph
COVER     := $(GOTOOLDIR)/cover
DIGRAPH   := $(GOTOOLDIR)/digraph
EG        := $(GOTOOLDIR)/eg
GODEX     := $(GOTOOLDIR)/godex
GODOC     := $(GOTOOLDIR)/godoc
GOIMPORTS := $(GOTOOLDIR)/goimports
GOMVPKG   := $(GOTOOLDIR)/gomvpkg
GOTYPE    := $(GOTOOLDIR)/gotype
ORACLE    := $(GOTOOLDIR)/oracle
SSADUMP   := $(GOTOOLDIR)/ssadump
STRINGER  := $(GOTOOLDIR)/stringer
VET       := $(GOTOOLDIR)/vet
GOX       := $(GOBIN)/gox
LINT      := $(GOBIN)/lint
GODEP     := $(GOBIN)/godep

$(BENCHCMP)  : ; @go get -v golang.org/x/tools/cmd/benchcmp
$(CALLGRAPH) : ; @go get -v golang.org/x/tools/cmd/callgraph
$(COVER)     : ; @go get -v golang.org/x/tools/cmd/cover
$(DIGRAPH)   : ; @go get -v golang.org/x/tools/cmd/digraph
$(EG)        : ; @go get -v golang.org/x/tools/cmd/eg
$(GODEX)     : ; @go get -v golang.org/x/tools/cmd/godex
$(GODOC)     : ; @go get -v golang.org/x/tools/cmd/godoc
$(GOIMPORTS) : ; @go get -v golang.org/x/tools/cmd/goimports
$(GOMVPKG)   : ; @go get -v golang.org/x/tools/cmd/gomvpkgs
$(GOTYPE)    : ; @go get -v golang.org/x/tools/cmd/gotype
$(ORACLE)    : ; @go get -v golang.org/x/tools/cmd/oracle
$(SSADUMP)   : ; @go get -v golang.org/x/tools/cmd/ssadump
$(STRINGER)  : ; @go get -v golang.org/x/tools/cmd/stringer
$(VET)       : ; @go get -v golang.org/x/tools/cmd/vet
$(LINT)      : ; @go get -v github.com/golang/lint/golint
$(GOX)       : ; @go get -v github.com/mitchellh/gox
$(GODEP)     : ; @go get -v github.com/tools/godep

.PHONY: vet
vet: $(VET) ; @for src in $(GOSOURCES); do GOPATH=$(ORIGGOPATH) go tool vet $$src; done

.PHONY: lint
lint: $(LINT) ; @for src in $(GOSOURCES); do GOPATH=$(ORIGGOPATH) golint $$src || exit 1; done

.PHONY: fmt
fmt: ; @GOPATH=$(ORIGGOPATH) go fmt ./...

.PHONY: test
test: ; @go test -v ./...

.PHONY: race
race: ; @for pkg in $(GOPKGS); do go test -v -race $$pkg || exit 1; done

.PHONY: deps
deps: ; @GOPATH=$(ORIGGOPATH) go get -u -v -t ./...

.PHONY: cover
cover: $(COVER)
	@echo 'mode: set' > $(COVERAGEOUT)
	@for pkg in $(GOPKGS); do \
		go test -v -coverprofile=$(COVERAGETMP) $$pkg || exit 1; \
		grep -v 'mode: set' $(COVERAGETMP) >> $(COVERAGEOUT); \
		rm $(COVERAGETMP); \
	done
	@go tool cover -html=$(COVERAGEOUT)

##########################################################################################
## Godep support
##########################################################################################

.PHONY: savegodeps
savegodeps: $(GODEP) ; @GOPATH=$(ORIGGOPATH) $(GODEP) save ./...

.PHONY: restoregodeps
restoregodeps: $(GODEP) ; @GOPATH=$(ORIGGOPATH) $(GODEP) restore

.PHONY: updategodeps
updategodeps: $(GODEP) ; @GOPATH=$(ORIGGOPATH) $(GODEP) update ./...

##########################################################################################
## Make utilities
##########################################################################################

.PHONY: printvars
printvars:
	@$(foreach V, $(sort $(.VARIABLES)), $(if $(filter-out environment% default automatic, $(origin $V)), $(warning $V=$($V) )))
	@exit 0
