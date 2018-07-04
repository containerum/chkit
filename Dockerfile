FROM golang:1.10-alpine AS builder
RUN apk --no-cache add zip make git openssl
WORKDIR /go/src/github.com/containerum/chkit

ARG BUILD_CONTAINERUM_API=https://api.containerum.io
ENV CONTAINERUM_API=$BUILD_CONTAINERUM_API \
    PATH=$PATH:$GOPATH/bin

COPY . .
RUN go get -u -v github.com/UnnoTed/fileb0x && \
    cd $GOPATH/src/github.com/UnnoTed/fileb0x && \
    git checkout 033c2ecc1c0f93d04afe94186f15193dd4441646 && \
    go install  && \
    cd $GOPATH/src/github.com/containerum/chkit

RUN make build

FROM alpine:3.7
RUN apk --no-cache add ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/containerum/chkit/build/chkit /chkit

VOLUME /root/.config/containerum

ENTRYPOINT ["/chkit"]
CMD ["--help"]
