FROM golang:latest

RUN apt-get update
#RUN apt-get upgrade -y

RUN apt-get install -y wget cmake

ADD . /go/src/github.com/welaw/welaw
WORKDIR /go/src/github.com/welaw/welaw

#RUN make install-deps
RUN make install-libgit2
RUN make install-go-deps
RUN make start

EXPOSE 8080

CMD ["welaw"]
