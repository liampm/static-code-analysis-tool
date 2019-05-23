FROM golang:1.12.5-alpine3.9

ENV GO111MODULE=on

RUN apk add git inotify-tools procps

ADD ./src /go/src/app
WORKDIR /go/src/app

EXPOSE ${PORT}

COPY ./go-watch.sh /bin/go-watch.sh

RUN chmod +x /bin/go-watch.sh

CMD ['go-watch.sh']


