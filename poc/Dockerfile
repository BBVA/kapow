FROM python:3.7-alpine
RUN apk update && apk add bash curl coreutils file
RUN pip install pipenv
COPY Pipfile Pipfile.lock /tmp/
RUN cd /tmp && pipenv install --system --deploy
COPY kapow /usr/bin
COPY bin/* /usr/bin/
ENTRYPOINT ["/usr/bin/kapow"]
