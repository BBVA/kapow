FROM python:3.7.3-alpine3.8

# Install Dependences
RUN pip install pipenv
COPY Pipfile Pipfile.lock ./
RUN pipenv install --system --deploy && rm Pipfile Pipfile.lock

# Install kapow!
COPY kapow.py /usr/sbin/kapow
RUN chmod a+x /usr/sbin/kapow
