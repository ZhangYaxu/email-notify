FROM golang:1.5.1

EXPOSE 8000

ADD . /go

RUN go install piccolo

ENTRYPOINT /go/bin/piccolo
