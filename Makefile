BINARY = chkit
PACKAGE = chkit-v2
COMMIT_HASH = `git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE = `date +%FT%T%Z`
DEFAULT_TCP_SERVER = sdk.containerum.io:3000
DEFAULT_HTTP_SERVER = http://sdk.containerum.io:3333
VERSION = 2.0.3
REQLDFLAGS = -X ${PACKAGE}/chlib.CommitHash=${COMMIT_HASH} \
	-X ${PACKAGE}/chlib.BuildDate=${BUILD_DATE} \
	-X ${PACKAGE}/chlib/dbconfig.DefaultTCPServer=${DEFAULT_TCP_SERVER} \
	-X ${PACKAGE}/chlib/dbconfig.DefaultHTTPServer=${DEFAULT_HTTP_SERVER} \
	-X ${PACKAGE}/chlib.DevGoPath=${GOPATH} \
	-X ${PACKAGE}/chlib.DevGoRoot=${GOROOT} \
	-X ${PACKAGE}/helpers.CurrentClientVersion=${VERSION}

BUILDDIR = build
#track sources
SOURCES = $(shell find ${PWD} -name '*.go')

define do_build
@echo -e "\x1b[35mRun go build\x1b[0m"
@go build -ldflags "${LDFLAGS} ${REQLDFLAGS}" -o ${1}
endef

#remove source file after packing
%.tar.gz : ${SOURCES}
	$(call do_build,${BUILDDIR}/${BINARY})
	@echo -e "\x1b[35mPack to $@\x1b[0m"
	@chmod +x ${BUILDDIR}/chkit
	@tar --transform 's/.*\///g' --remove-files -cvzf $@ ${BUILDDIR}/${BINARY}

#removes source file after packing
%.zip : ${SOURCES}
	$(call do_build,${BUILDDIR}/${BINARY}.exe)
	@echo -e "\x1b[35mPack to $@\x1b[0m"
	@zip -jmD $@ ${BUILDDIR}/${BINARY}.exe

all: build

#for debugging purposes
build:
	go build -ldflags "${LDFLAGS} ${REQLDFLAGS}" -o ${BINARY}

clean:
	@if [ -f ${BINARY} ]; then rm ${BINARY}; fi

clean-build:
	@rm -rf ${BUILDDIR}

clean-all: clean clean-build

test:
	@go test

define custom_os_arch_build
	$(eval GOOS=${1})
	$(eval GOARCH=${2})
	@export GOOS GOARCH
	$(eval TARGET=${BINARY}_${GOOS}_${GOARCH}_v${VERSION})
	$(if $(filter ${GOOS},windows),$(eval TARGET=${TARGET}.zip),$(eval TARGET=${TARGET}.tar.gz))
	$(eval TARGET=$(subst darwin,mac,${TARGET}))
	$(eval TARGET=$(subst 386,x86,${TARGET}))
	$(eval TARGET=$(subst amd64,x64,${TARGET}))
	@echo -e "\x1b[32;1mBuild ${TARGET}\x1b[0m"
	@$(MAKE) -s -f $(lastword $(MAKEFILE_LIST)) LDFLAGS="-w -s" ${BUILDDIR}/${TARGET}

endef

#production builds
build-all:
	@mkdir -p build
	$(call custom_os_arch_build,linux,386)
	$(call custom_os_arch_build,linux,amd64)
	$(call custom_os_arch_build,linux,arm)
	$(call custom_os_arch_build,darwin,amd64)
	$(call custom_os_arch_build,windows,386)
	$(call custom_os_arch_build,windows,amd64)

.PHONY: all build clean clean-build clean-all