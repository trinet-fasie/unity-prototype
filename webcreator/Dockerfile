FROM node:8.12-alpine
MAINTAINER TrinetMarketing

RUN apk --update add \
    git \
    gcc \
    g++ \
    make \
    python \
    openssh && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

COPY ./frontend /var/www/webcreator
RUN rm -rf /var/www/webcreator/node_modules
RUN chmod +x /var/www/webcreator/wait-for.sh

WORKDIR /var/www/webcreator

RUN npm install

EXPOSE 8080

CMD \
    npm install && \
    npm run serve