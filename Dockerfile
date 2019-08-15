FROM golang:alpine as builder
RUN apk add -U --no-cache ca-certificates
RUN mkdir /go/src/drone-rancher
ADD . /go/src/drone-rancher/
WORKDIR /go/src/drone-rancher
RUN apk add git
RUN go get -d -v
RUN rm -rf /go/src/github.com/rancher/types/vendor/github.com/rancher
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o drone-rancher .

FROM scratch

ENV GODEBUG=netdns=go

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

LABEL org.label-schema.version=latest
LABEL org.label-schema.vcs-url="https://github.com/JTPearce/drone-rancher.git"
LABEL org.label-schema.name="Drone Rancher"
LABEL org.label-schema.vendor="James Pearce"

COPY --from=builder /go/src/drone-rancher/drone-rancher /
ENTRYPOINT ["/drone-rancher"]
