FROM ubuntu:14.04
MAINTAINER Dmitry Vorobev <dimahabr@gmail.com>

# update packages, install dependencies and set GOPATH
RUN apt-get update -q
RUN apt-get install -qy curl libvips-dev git rabbitmq-server
RUN curl -s https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz | tar -v -C /usr/local -xz
ENV PATH /usr/local/go/bin:/go/bin:$PATH
ENV GOPATH /go

# copy application to container
ADD . /go/src/github.com/pydima/go-thumbnailer
WORKDIR /go/src/github.com/pydima/go-thumbnailer
RUN go get
RUN mkdir -p /etc/go_thumbnailer && cp /go/src/github.com/pydima/go-thumbnailer/config.json /etc/go_thumbnailer/config.json
ENTRYPOINT /bin/bash
