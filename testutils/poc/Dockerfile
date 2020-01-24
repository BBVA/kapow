FROM python:3.7-alpine

RUN apk upgrade --update-cache       \
 && apk add bash curl coreutils file \
 && pip install pipenv

COPY Pipfile Pipfile.lock /

RUN pipenv install --system --deploy \
 && rm -f /Pipfile /Pipfile.lock

COPY bin/* /usr/bin/

ENTRYPOINT ["/usr/bin/kapow"]
