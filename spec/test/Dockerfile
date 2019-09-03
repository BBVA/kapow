FROM python:3.7-alpine

# Install CircleCI requirements for base images
#   https://circleci.com/docs/2.0/custom-images/
RUN apk upgrade --update-cache \
 && apk add git openssh-server tar gzip ca-certificates

# Install Kapow! Spec Test Suite
RUN mkdir -p /usr/src/ksts
WORKDIR /usr/src/ksts
COPY features /usr/src/ksts/features
COPY Pipfile Pipfile.lock /usr/src/ksts/
RUN pip install --upgrade pip \
 && pip install pipenv \
 && pipenv install --deploy --system \
 && rm -f Pipfile Pipfile.lock
