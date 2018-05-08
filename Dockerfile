FROM golang:1.10-alpine AS builder
RUN apk --no-cache add zip make git openssl
WORKDIR /go/src/github.com/containerum/chkit
ARG BUILD_CONTAINERUM_API=https://api.containerum.io:8082
ENV CONTAINERUM_API=$BUILD_CONTAINERUM_API
COPY . .
RUN make build

FROM alpine:3.7
RUN apk --no-cache add ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/containerum/chkit/build/chkit /chkit
VOLUME /root/.config/containerum
ARG ALLOW_SELF_SIGNED_CERTS=false
ARG USERNAME
ARG PASS
RUN echo $ALLOW_SELF_SIGNED_CERTS
RUN /chkit set containerum-api --allow-self-signed-certs=$ALLOW_SELF_SIGNED_CERTS
RUN /chkit login -u $USERNAME -p $PASS --default-namespace -
ENTRYPOINT ["/chkit"]