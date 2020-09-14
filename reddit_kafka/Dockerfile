FROM ubuntu

# Install the C lib for kafka
RUN apt-get update
RUN apt-get install -y --no-install-recommends apt-utils wget gnupg software-properties-common
RUN apt-get install -y apt-transport-https ca-certificates
RUN wget -qO - https://packages.confluent.io/deb/5.1/archive.key | apt-key add -
RUN add-apt-repository "deb [arch=amd64] https://packages.confluent.io/deb/5.1 stable main"
RUN apt-get update
RUN apt-get install -y librdkafka-dev

# Install Go
RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get update
RUN apt-get install -y golang-1.14-go


# build the library
WORKDIR /go/src/gitlab.appsflyer.com/rantav/kafka-mirror-tester
#COPY *.go ./
COPY . .
#COPY vendor vendor

RUN GOPATH=/go GOOS=linux /usr/lib/go-1.14/bin/go build -a -o reddit-kafka .

#EXPOSE 9092

ENTRYPOINT ["./reddit-kafka"]
#
#FROM golang:1.14-alpine as builder
#
#WORKDIR /reddit_kafka
#
#COPY . .
#
#RUN go get -u all
#RUN go build -o reddit-kafka
#
#
#FROM alpine:latest
#
#WORKDIR /reddit_kafka
#
#COPY --from=builder /reddit_kafka/reddit-kafka ./reddit-kafka
#
#EXPOSE 8010
#
#CMD ["reddit-kafka"]
