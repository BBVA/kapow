FROM bbva/kapow:0.1
RUN apk update && apk add nmap
COPY nmap.pow /tmp/
RUN cd /tmp && pipenv install --system --deploy
ENTRYPOINT ["/usr/bin/kapow"]
CMD ["server", "/tmp/nmap.pow"]
