FROM golang:1.10-alpine AS builder
RUN apk --no-cache add zip make ca-certificates git openssl && update-ca-certificates
WORKDIR /go/src/github.com/containerum/chkit
ENV CONTAINERUM_API https://api.containerum.io:8082
COPY . .
RUN make build

FROM alpine:3.7
COPY --from=builder /usr/local/share/ca-certificates /usr/local/share/ca-certificates
COPY --from=builder /go/src/github.com/containerum/chkit/build/chkit /chkit
VOLUME /root/.config/containerum
ENTRYPOINT ["/chkit"]