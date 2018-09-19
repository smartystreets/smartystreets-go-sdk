FROM golang:alpine

ADD . /code
WORKDIR /code

RUN apk add --upgrade make git gcc musl-dev \
	&& wget -O - "https://github.com/smartystreets/version-tools/releases/download/0.0.6/release.tar.gz" | tar -xz -C /usr/local/bin/
