FROM alpine

RUN apk upgrade --update-cache       \
    && apk add bash curl coreutils file

COPY kapow /usr/bin/kapow

COPY *.pow /tmp/

ENTRYPOINT ["/usr/bin/kapow"]
