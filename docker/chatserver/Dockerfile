FROM golang:alpine 
MAINTAINER Louis DeLosSantos

COPY . /go/src/github.com/ldelossa/rtchat

# Confirm vendor folder exists, expected to run dependency management tools outside of container
RUN ls /go/src/github.com/ldelossa/rtchat/vendor || exit 1

RUN cd /go/src/github.com/ldelossa/rtchat/cmd/chatserver && \
	go install

ENTRYPOINT ["/go/bin/chatserver"]
EXPOSE 8080
