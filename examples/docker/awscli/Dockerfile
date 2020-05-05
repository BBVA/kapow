FROM bbvalabsci/kapow:latest as kp
FROM amazon/aws-cli:latest

COPY --from=kp /kapow /usr/bin/kapow

ENTRYPOINT kapow
