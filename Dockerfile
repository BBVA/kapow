FROM golang:1.13 as build

RUN go get github.com/BBVA/kapow
RUN CGO_ENABLED=0 GOOS=linux go install github.com/BBVA/kapow

FROM alpine:latest
COPY --from=build /go/bin/kapow /usr/bin/kapow

ENTRYPOINT /usr/bin/kapow
