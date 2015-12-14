FROM golang:1

ADD . /gorepos
WORKDIR /gorepos

RUN go get -d .
RUN go build ./...

CMD /gorepos/gorepos
