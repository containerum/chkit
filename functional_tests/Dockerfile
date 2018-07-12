FROM python:3.6 as tester
RUN apt-get update && \
    apt-get install -y sudo snap apt-utils
RUN wget -O go.tar.gz https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xvf go.tar.gz && \
    rm go.tar.gz
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH
RUN go get -u -v -d github.com/UnnoTed/fileb0x && \
        cd $GOPATH/src/github.com/UnnoTed/fileb0x && \
        git checkout 033c2ecc1c0f93d04afe94186f15193dd4441646 && \
        go install
RUN go get -u -v -d github.com/containerum/chkit/cmd/chkit
WORKDIR $GOPATH/src/github.com/containerum/chkit

ENV CONTAINERUM_API="http://local.dev"
ENV TEST_USER=""
ENV TEST_USER_PASSWORD=""
ENV TEST_NAMESPACE=""

RUN git checkout chkit-v3 && \
    make genkey && \
    make build && \
    printf "%s test-host.hub.containerum.io\n" $CONTAINERUM_API >> /etc/hosts && \
    pip3 install -r functional_tests/requirements.txt
CMD bash -c "alias python=python3 && make functional_tests"