FROM postgres:10-alpine
MAINTAINER TrinetMarketing

RUN apk --update add bash && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY ./migrations /migrations

COPY go-migrate /usr/local/bin/go-migrate
RUN chmod +x /usr/local/bin/go-migrate
COPY migrate.sh /usr/local/bin/migrate
RUN chmod +x /usr/local/bin/migrate
COPY create.sh /usr/local/bin/create
RUN chmod +x /usr/local/bin/create
