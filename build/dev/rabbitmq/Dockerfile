FROM rabbitmq:3-management-alpine
MAINTAINER TrinetMarketing

COPY create-user.sh /usr/local/bin/create-user.sh
RUN chmod +x /usr/local/bin/create-user.sh

RUN rabbitmq-plugins enable rabbitmq_web_stomp

