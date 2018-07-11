.PHONY: docker genkey build test clean install release single_release dev mock functional_tests

CMD_DIR:=cmd/chkit
CLI_DIR:=pkg/cli
#get current package, assuming it`s in GOPATH sources
PACKAGE := $(shell go list -f '{{.ImportPath}}' ./$(CLI_DIR))
PACKAGE := $(PACKAGE:%/$(CLI_DIR)=%)
SIGNING_KEY_DIR:=~/.config/containerum/.chkit-sign
PRIVATE_KEY_FILE:=privkey.pem
PUBLIC_KEY_FILE:=pubkey.pem
FILEBOX := $(shell command -v fileb0x 2>/dev/null)
FUNCTIONAL_TEST_MODULES := config deployment pod service configmap solution
HELP_DIR := help
HELP_CONTENT_FILES := $(shell find $(HELP_DIR)/content -name '*.md')
PIP := pip3
PYTHON := python3

COMMIT_HASH=$(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE=$(shell date +%FT%T%Z)
LATEST_TAG=$(shell git describe --tags $(shell git rev-list --tags --max-count=1))
CONTAINERUM_API?=https://api.containerum.io
VERSION?=$(LATEST_TAG:v%=%)

# make directory and store path to variable
BUILDS_DIR:=$(PWD)/build
EXECUTABLE:=chkit
RAW_PUBLIC_KEY:=$(shell openssl enc -base64 -in $(SIGNING_KEY_DIR)/$(PUBLIC_KEY_FILE))
SPACE:=$(shell echo ' ')
PUBLIC_KEY:=$(subst $(SPACE),,$(RAW_PUBLIC_KEY))
RELEASE_LDFLAGS=-X $(PACKAGE)/$(CLI_DIR).VERSION=v$(VERSION) \
	-X $(PACKAGE)/pkg/update.PublicKeyB64=$(PUBLIC_KEY) \
	-X $(PACKAGE)/$(CLI_DIR)/mode.API_ADDR=$(CONTAINERUM_API)
DEV_LDFLAGS=-X '$(PACKAGE)/$(CLI_DIR)/mode.API_ADDR=$(CONTAINERUM_API)' \
	-X '$(PACKAGE)/$(CLI_DIR).VERSION=v$(VERSION)' \
	-X $(PACKAGE)/pkg/update.PublicKeyB64=$(PUBLIC_KEY)

CONTAINER_NAME?=containerum/chkit
ALLOW_SELF_SIGNED_CERTS?=true
docker:
	docker build -t $(CONTAINER_NAME) . \
		--build-arg ALLOW_SELF_SIGNED_CERTS=$(ALLOW_SELF_SIGNED_CERTS)

$(SIGNING_KEY_DIR)/$(PRIVATE_KEY_FILE):
	@echo "Generating private/public ECDSA keys to sign"
	@mkdir -p $(SIGNING_KEY_DIR)
	@openssl ecparam -genkey -name prime256v1 -out $(SIGNING_KEY_DIR)/$(PRIVATE_KEY_FILE)
	@openssl ec -in $(SIGNING_KEY_DIR)/$(PRIVATE_KEY_FILE) -pubout -out $(SIGNING_KEY_DIR)/$(PUBLIC_KEY_FILE)
	@echo "Keys stored in $(SIGNING_KEY_DIR)"

genkey: $(SIGNING_KEY_DIR)/$(PRIVATE_KEY_FILE)

help/ab0x.go: help/b0x.toml $(HELP_CONTENT_FILES)
ifndef FILEBOX
	$(error "fileb0x is not available, please install it from https://github.com/UnnoTed/fileb0x)
endif
	go generate ./help

# go has build artifacts caching so soruce tracking not needed
build: $(SIGNING_KEY_DIR)/$(PRIVATE_KEY_FILE) help/ab0x.go
	@echo "Building chkit for current OS/architecture, without signing"
	go build -v -ldflags="$(RELEASE_LDFLAGS)" -o $(BUILDS_DIR)/$(EXECUTABLE) ./$(CMD_DIR)

test:
	@echo "Running tests"
	@go test -v ./...

clean:
	@rm -rf $(BUILDS_DIR)

install: help/ab0x.go
	@go install -v -ldflags="$(RELEASE_LDFLAGS)" ./$(CMD_DIR)

# lambda to generate build dir name using os,arch,version
temp_dir_name=$(EXECUTABLE)_$(1)_$(2)_v$(3)

pack_win=zip -r -j $(1).zip $(1) && rm -rf $(1)
# pack_nix=tar --transform 's/.*\///g' -cpzf $(1).tar.gz $(1)/* && rm -rf $(1)
pack_nix=tar -C $(1) -cpzf $(1).tar.gz ./ && rm -rf $(1)
create_checksum=openssl dgst -sha256 -binary -out $(dir $(1))/sha256.sum $(1)
create_signature=openssl dgst -sha256 -sign $(SIGNING_KEY_DIR)/$(PRIVATE_KEY_FILE) -out $(dir $(1))/ecdsa.sig $(1)

define build_release
@echo "Building release package for OS $(1), arch $(2)"
$(eval temp_build_dir=$(BUILDS_DIR)/$(call temp_dir_name,$(1),$(2),$(VERSION)))
@mkdir -p $(temp_build_dir)
$(eval ifeq ($(1),windows)
	temp_executable=$(temp_build_dir)/$(EXECUTABLE).exe
else
	temp_executable=$(temp_build_dir)/$(EXECUTABLE)
endif)
GOOS=$(1) GOARCH=$(2) go build -tags="release" -ldflags="$(RELEASE_LDFLAGS)"  -v -o $(temp_executable) ./$(CMD_DIR)
@$(call create_checksum,$(temp_executable))
@$(call create_signature,$(temp_executable))
$(eval ifeq ($(1),windows)
	pack_cmd = $(call pack_win,$(temp_build_dir))
else
	pack_cmd = $(call pack_nix,$(temp_build_dir))
endif)
@$(pack_cmd)
endef

release: $(SIGNING_KEY_DIR)/$(PRIVATE_KEY_FILE) help/ab0x.go
	$(call build_release,linux,amd64)
	$(call build_release,linux,386)
	$(call build_release,linux,arm)
	$(call build_release,darwin,amd64)
	$(call build_release,windows,amd64)
	$(call build_release,windows,386)

single_release: $(SIGNING_KEY_DIR)/$(PRIVATE_KEY_FILE) help/ab0x.go
	$(call build_release,$(OS),$(ARCH))

dev: help/ab0x.go
	$(eval VERSION=$(LATEST_TAG:v%=%)+dev)
	@echo building $(VERSION)
	go build -v -race --tags="dev" --ldflags="$(DEV_LDFLAGS)" ./$(CMD_DIR)

mock: help/ab0x.go
	$(eval VERSION=$(LATEST_TAG:v%=%)+mock)
	@echo building $(VERSION)
	@go build -v --tags="dev mock" -ldflags="$(DEV_LDFLAGS)" ./$(CMD_DIR)

functional_tests: install
	@$(PIP) install -r functional_tests/requirements.txt
	@$(PYTHON) -m unittest $(foreach module,$(FUNCTIONAL_TEST_MODULES),functional_tests.$(module) ) -v
