FROM golang:1.11
LABEL maintainer="dkargashin3@gmail.com"

ENV GOPATH /go
ENV PATH ${GOPATH}/bin:$PATH
RUN go get -u github.com/golang/lint/golint

RUN apt-get update && apt-get install -y --no-install-recommends \
    clang \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN apt-get install -y sudo \
    ls -al /bin/sh && sudo rm /bin/sh \
    && sudo ln -s /bin/bash /bin/sh \
    && ls -al /bin/sh