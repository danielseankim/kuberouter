FROM golang:1.7

ADD . /go/src/github.com/logankimmel/oxy_dynamic_router

RUN cd /go/src/github.com/logankimmel/oxy_dynamic_router &&  go get ./...
RUN go install github.com/logankimmel/oxy_dynamic_router
WORKDIR /go/src/github.com/logankimmel/oxy_dynamic_router
ENTRYPOINT /go/bin/oxy_dynamic_router

