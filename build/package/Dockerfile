FROM golang

ADD . /go/src/github.com/a1k24/short-url/

RUN go get github.com/a1k24/short-url/cmd/server
RUN go install github.com/a1k24/short-url/cmd/server

ARG USERNAME
ARG PASSWORD
ARG DOMAIN
ENV USERNAME=${USERNAME}
ENV PASSWORD=${PASSWORD}
ENV DOMAIN=${DOMAIN}
# Run the server command by default when the container starts.
ENTRYPOINT /go/bin/server -username $USERNAME -password $PASSWORD -domain $DOMAIN

# Document that the service listens on port 10000.
EXPOSE 10000