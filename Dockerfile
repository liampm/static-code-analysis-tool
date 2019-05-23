FROM golang:1.12.5-alpine3.9

RUN apk add git

ADD ./src /go/src/app
WORKDIR /go/src/app

EXPOSE ${PORT}

CMD ["go", "run", "main.go"]


