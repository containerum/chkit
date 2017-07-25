BINARY = chkit
PACKAGE = chkit-v2
COMMIT_HASH = `git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE = `date +%FT%T%Z`
DEFAULT_TCP_SERVER = sdk.containerum.io:3000
DEFAULT_HTTP_SERVER = http://sdk.containerum.io:3333
VERSION = "2.0.1"
LDFLAGS = "-X ${PACKAGE}/chlib.CommitHash=${COMMIT_HASH} \
	-X ${PACKAGE}/chlib.BuildDate=${BUILD_DATE} \
	-X ${PACKAGE}/chlib/dbconfig.DefaultTCPServer=${DEFAULT_TCP_SERVER} \
	-X ${PACKAGE}/chlib/dbconfig.DefaultHTTPServer=${DEFAULT_HTTP_SERVER} \
	-X ${PACKAGE}/chlib.DevGoPath=${GOPATH} \
	-X ${PACKAGE}/chlib.DevGoRoot=${GOROOT} \
	-X ${PACKAGE}/helpers.CurrentClientVersion=${VERSION}"

all: build

build:
	go build -ldflags ${LDFLAGS} -o ${BINARY}

clean:
	if [ -f ${BINARY} ]; then rm ${BINARY}; fi
test:
	go test

build-all: clean-build
	mkdir build
	#Build for linux_386
	GOOS=linux GOARCH=386 go build -o build/chkit_linux_x86_v${VERSION} -ldflags ${LDFLAGS}
	#Build for linux_amd64
	GOOS=linux GOARCH=amd64 go build -o build/chkit_linux_x64_v${VERSION} -ldflags ${LDFLAGS}
	#Build for darwin_amd64
	GOOS=darwin GOARCH=amd64 go build -o build/chkit_mac_x64_v${VERSION} -ldflags ${LDFLAGS}
	#Build for windows_386
	GOOS=windows GOARCH=386 go build -o build/chkit_win_x86_v${VERSION}.exe -ldflags ${LDFLAGS}
	#Build for windows_amd64
	GOOS=windows GOARCH=amd64 go build -o build/chkit_win_x64_v${VERSION}.exe -ldflags ${LDFLAGS}
	#Build for linux ARM
	GOOS=linux GOARCH=arm go build -o build/chkit_linux_arm_v${VERSION} -ldflags ${LDFLAGS}
clean-build:
	rm -rf build

.PHONY: all build clean
