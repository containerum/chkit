BINARY = chkit
PACKAGE = github.com/containerum/chkit
COMMIT_HASH = $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE = $(shell date +%FT%T%Z)
DEFAULT_TCP_SERVER = sdk.containerum.io:3000
DEFAULT_HTTP_SERVER = http://sdk.containerum.io:3333
VERSION = $(shell git describe --tags $(shell git rev-list --tags --max-count=1))
REQLDFLAGS = -X ${PACKAGE}/chlib.CommitHash=${COMMIT_HASH} \
	-X ${PACKAGE}/chlib.BuildDate=${BUILD_DATE} \
	-X ${PACKAGE}/chlib/dbconfig.DefaultTCPServer=${DEFAULT_TCP_SERVER} \
	-X ${PACKAGE}/chlib/dbconfig.DefaultHTTPServer=${DEFAULT_HTTP_SERVER} \
	-X ${PACKAGE}/helpers.CurrentClientVersion=$(VERSION:v%=%)

BUILDDIR = build
USE_DOCKER ?= 1
#track sources
SOURCES = $(shell find ${PWD} -name '*.go')

#for installation
PREFIX ?= usr
DESTDIR ?=
INSTDIR ?= ${DESTDIR}/${PREFIX}/bin
AUTOCOMPDIR ?= ${DESTDIR}/${PREFIX}/share/bash-completion/completions
AUTOCOMPFILE = ${AUTOCOMPDIR}/chkit.completion

ifeq ($(USE_DOCKER),1)
have_docker_perm = $(shell docker images >/dev/null 2>/dev/null && echo 1)
ifneq ($(have_docker_perm),1)
$(info You don`t have permissions to run docker, so run with sudo)
sudo=sudo
endif
define do_build
@echo -e "\x1b[35mRun go build\x1b[0m"
@$(sudo) docker run --rm \
	-v $(shell pwd):/go/src/${PACKAGE} \
	-w /go/src/${PACKAGE} \
	-e GOOS=${GOOS} \
	-e GOARCH=${GOARCH} \
	-it golang:1.9 \
	/bin/bash -c "go build -v -ldflags '${LDFLAGS} ${REQLDFLAGS}' -o ${1} && \
		chown $(shell id -u) ${1}"
endef
else
define do_build
@echo -e "\x1b[35mRun go build\x1b[0m"
go build -v -ldflags '${LDFLAGS} ${REQLDFLAGS}' -o ${1}
endef
endif


${BUILDDIR}/${BINARY}: ${SOURCES}
	$(call do_build,${BUILDDIR}/${BINARY})

#remove source file after packing
%.tar.gz : ${SOURCES}
	$(call do_build,${BUILDDIR}/${BINARY})
	@echo -e "\x1b[35mPack to $@\x1b[0m"
	@chmod +x ${BUILDDIR}/chkit
	@tar --transform 's/.*\///g' --remove-files -cpvzf $@ ${BUILDDIR}/${BINARY}

#removes source file after packing
%.zip : ${SOURCES}
	$(call do_build,${BUILDDIR}/${BINARY}.exe)
	@echo -e "\x1b[35mPack to $@\x1b[0m"
	@zip -jmD $@ ${BUILDDIR}/${BINARY}.exe

all: ${BUILDDIR}/${BINARY}

clean:
	@rm -rf ${BUILDDIR}

test:
	@go test

define custom_os_arch_build
	$(eval GOOS=${1})
	$(eval GOARCH=${2})
	$(eval TARGET=${BINARY}_${GOOS}_${GOARCH}_v${VERSION})
	$(if $(filter ${GOOS},windows),$(eval TARGET=${TARGET}.zip),$(eval TARGET=${TARGET}.tar.gz))
	$(eval TARGET=$(subst darwin,mac,${TARGET}))
	$(eval TARGET=$(subst 386,x86,${TARGET}))
	$(eval TARGET=$(subst amd64,x64,${TARGET}))
	@echo -e "\x1b[32;1mBuild ${TARGET}\x1b[0m"
	@$(MAKE) -s -f $(lastword $(MAKEFILE_LIST)) GOOS=${GOOS} GOARCH=${GOARCH} LDFLAGS="-w -s" ${BUILDDIR}/${TARGET}

endef

#production builds
build-all:
	@mkdir -p ${BUILDDIR}
	$(call custom_os_arch_build,linux,386)
	$(call custom_os_arch_build,linux,amd64)
	$(call custom_os_arch_build,linux,arm)
	$(call custom_os_arch_build,darwin,amd64)
	$(call custom_os_arch_build,windows,386)
	$(call custom_os_arch_build,windows,amd64)

install: ${BUILDDIR}/${BINARY}
	@echo -e "\x1b[32;1mInstall $< to ${INSTDIR}\x1b[0m"
	@echo -e "\x1b[35mCopy $< to ${INSTDIR}\x1b[0m"
	@install -D $< ${INSTDIR}/${BINARY}
	@echo -e "\x1b[35mGenerate autocomplete\x1b[0m"
	@install -d ${AUTOCOMPDIR}
	@$< genautocomplete -f ${AUTOCOMPFILE}

uninstall:
	@echo -e "\x1b[32;1mUninstall $< from ${INSTDIR}\x1b[0m"
	@rm ${INSTDIR}/${BINARY}
	@echo -e "\x1b[35mRemove autocomplete ${AUTOCOMPFILE}\x1b[0m"
	@rm ${AUTOCOMPFILE}
	@echo -e "\x1b[35mNo removing config dir ${HOME}/.containerum\x1b[0m"

.PHONY: all build clean clean-build clean-all
