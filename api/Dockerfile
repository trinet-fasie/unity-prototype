FROM golang:1.11-alpine
MAINTAINER TrinetMarketing

RUN apk --update add git openssh && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

COPY . /go/src/github.com/trinet-fasie/unity-prototype/api
RUN chmod +x /go/src/github.com/trinet-fasie/unity-prototype/api/wait-for.sh

WORKDIR /go/src/github.com/trinet-fasie/unity-prototype/api

RUN go get -d -v && go install -v

EXPOSE 3000

CMD \
    /go/src/github.com/trinet-fasie/unity-prototype/api/wait-for.sh db:5432 --timeout=60 && \
    /go/src/github.com/trinet-fasie/unity-prototype/api/wait-for.sh rabbitmq:5672 --timeout=60 && \
    api
