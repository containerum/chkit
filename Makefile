.PHONY: genkey genkey build test clean release single_release

#get current package, assuming it`s in GOPATH sources
PACKAGE := $(PWD:$(GOPATH)/src/%=%)

SIGNING_KEY_DIR:=~/.chkit-sign
PRIVATE_KEY_FILE:=privkey.pem
PUBLIC_KEY_FILE:=pubkey.pem

COMMIT_HASH=$(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE=$(shell date +%FT%T%Z)
LATEST_TAG=$(shell git describe --tags $(shell git rev-list --tags --max-count=1))

VERSION?=$(LATEST_TAG:v%=%)

# make directory and store path to variable
BUILDS_DIR:=$(PWD)/build
EXECUTABLE:=chkit
LDFLAGS=-X $(PACKAGE)/cmd.Version=$(VERSION) \
	-X $(PACKAGE)/pkg/update.PublicKeyB64=\'$(shell base64 -w 0 $(SIGNING_KEY_DIR)/$(PUBLIC_KEY_FILE))\'

genkey:
	@echo "Generating private/public ECDSA keys to sign"
	@mkdir -p $(SIGNING_KEY_DIR)
	@openssl ecparam -genkey -name prime256v1 -out $(SIGNING_KEY_DIR)/$(PRIVATE_KEY_FILE)
	@openssl ec -in $(TEMP) -pubout -out $(SIGNING_KEY_DIR)/$(PUBLIC_KEY_FILE)
	@echo "Keys stored in $(SIGNING_KEY_DIR)"

# go has build artifacts caching so soruce tracking not needed
build:
	@echo "Building chkit for current OS/architecture, without signing"
	@go build -v -ldflags="$(LDFLAGS)" -o $(BUILDS_DIR)/$(EXECUTABLE)

test:
	@echo "Running tests"
	@go test -v ./...

clean:
	@rm -rf $(BUILDS_DIR)

install:
	@go install -ldflags="$(LDFLAGS)"

# lambda to generate build dir name using os,arch,version
temp_dir_name=$(EXECUTABLE)_$(1)_$(2)_v$(3)

pack_win=zip -r -j $(1).zip $(1) && rm -rf $(1)
pack_nix=tar --transform 's/.*\///g' -cpzf $(1).tar.gz $(1)/* && rm -rf $(1)

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
@GOOS=$(1) GOARCH=$(2) go build -ldflags="$(LDFLAGS)" -v -o $(temp_executable)
@$(call create_checksum,$(temp_executable))
@$(call create_signature,$(temp_executable))
$(eval ifeq ($(1),windows)
	pack_cmd = $(call pack_win,$(temp_build_dir))
else
	pack_cmd = $(call pack_nix,$(temp_build_dir))
endif)
@$(pack_cmd)
endef

release:
	$(call build_release,linux,amd64)
	$(call build_release,linux,386)
	$(call build_release,linux,arm)
	$(call build_release,darwin,amd64)
	$(call build_release,windows,amd64)
	$(call build_release,windows,386)

single_release:
	$(call build_release,$(OS),$(ARCH))
