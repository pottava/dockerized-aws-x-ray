FROM pottava/golang:1.9-builder AS build-env
ADD . /go/src/github.com/pottava/dockerized-aws-x-ray/sample/
WORKDIR /go/src/github.com/pottava/dockerized-aws-x-ray/sample/src
RUN go install

FROM alpine:3.6
ENV AWS_XRAY_DAEMON_ADDRESS=xray:2000
RUN apk add --no-cache ca-certificates

COPY --from=build-env /go/bin/src /usr/local/bin/xray

EXPOSE 80
ENTRYPOINT ["xray"]
