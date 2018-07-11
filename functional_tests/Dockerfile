FROM containerum/chkit:latest as chkit
# ---
FROM python:3.7.0-slim-stretch as tester
COPY --from=chkit /chkit /chkit-bin/chkit
ENV PATH = $PATH:/chkit-bin
RUN apt-get update && \
    apt-get install -y git make && \
    git clone https://github.com/containerum/chkit.git ./chkit-repo && \
    cd ./chkit-repo/chkit && \
    git checkout chkit-v3
WORKDIR ./chkit-repo
RUN make functional_tests