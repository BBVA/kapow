FROM golang:1.12.5-alpine3.9
RUN apk add git
RUN go get github.com/yudai/gotty

FROM kapow:latest
RUN apk add tmux
COPY --from=0 /go/bin/gotty .

RUN mkdir assets
COPY index.html ./assets
COPY slides.html ./assets
COPY kapow.md ./assets
COPY sbackground.png ./assets
COPY logo.png ./assets

COPY tutorial.pow .

COPY profile /etc/profile
COPY tmux.conf /root/.tmux.conf

CMD ["server", "tutorial.pow"]
