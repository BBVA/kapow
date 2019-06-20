FROM python:3.7-alpine

COPY Pipfile Pipfile.lock /tmp/

COPY bin/* /usr/bin/

WORKDIR /tmp

RUN apk upgrade --update-cache;                 \
    apk add                                     \
        bash                                    \
        curl                                    \
        coreutils                               \
        file;                                   \
                                                \
    pip install pipenv;                         \
                                                \
    pipenv install --system --deploy;

ENTRYPOINT ["/usr/bin/kapow"]
