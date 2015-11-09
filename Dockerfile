FROM golang:1.5.1

EXPOSE 8000

COPY . /usr/src/app

RUN go get -d -v && go install -v

CMD ["app"]
