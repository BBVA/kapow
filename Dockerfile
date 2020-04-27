FROM golang:1.14 as build

RUN go get github.com/BBVA/kapow
RUN CGO_ENABLED=0 GOOS=linux go install github.com/BBVA/kapow

FROM scratch
COPY --from=build /go/bin/kapow /kapow
