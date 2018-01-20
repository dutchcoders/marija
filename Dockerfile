FROM golang:latest

ARG LDFLAGS=""
RUN mkdir /config/
ADD config-docker.toml /config/config.toml
RUN go build -ldflags="$LDFLAGS" -o /go/bin/app github.com/dutchcoders/marija

ADD . /go/src/github.com/dutchcoders/marija

RUN go build -o /go/bin/marija github.com/dutchcoders/marija

ENTRYPOINT ["/go/bin/marija", "-p", "0.0.0.0:8080", "--config", "/config/config.toml"]
EXPOSE 8080
