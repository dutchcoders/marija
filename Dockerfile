FROM golang:latest

ADD . /go/src/github.com/dutchcoders/marija

RUN go build -o /go/bin/marija github.com/dutchcoders/marija

ENTRYPOINT /go/bin/marija -p 0.0.0.0:8080

EXPOSE 8080
