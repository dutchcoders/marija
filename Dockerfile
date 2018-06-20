FROM golang:1.9.4 AS go

RUN apt update -y

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

ADD . /go/src/github.com/dutchcoders/marija

ARG LDFLAGS=""

WORKDIR /go/src/github.com/dutchcoders/marija
RUN go build -ldflags="$(go run scripts/gen-ldflags.go)" -o /go/bin/app github.com/dutchcoders/marija

FROM debian
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=go /go/bin/app /marija/marija

ARG LDFLAGS=""
RUN mkdir /config/
ADD config-docker.toml /config/config.toml

ENTRYPOINT ["/marija/marija", "-p", "0.0.0.0:8080", "--config", "/config/config.toml"]
EXPOSE 8080
