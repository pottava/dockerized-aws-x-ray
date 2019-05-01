# AWS X-Ray v3.0
#
# docker run --rm pottava/xray:3.0 --help
# docker run --name xray -d -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY \
#     -p 2000:2000/udp pottava/xray:3.0 --region ${AWS_REGION} --local-mode --log-level debug

FROM alpine:3.9 AS certs
RUN apk --no-cache add "ca-certificates=20190108-r0"

FROM golang:1.12.4-alpine3.9 AS build
RUN apk --no-cache add "git=2.20.1-r0"
RUN go get -u github.com/aws/aws-xray-daemon/...
WORKDIR /go/src/github.com/aws/aws-xray-daemon
ENV XRAY_VERSION=3.0.1
RUN git checkout "v${XRAY_VERSION}"
WORKDIR /go/src/github.com/aws/aws-xray-daemon/daemon
RUN go build -ldflags "-s -w" -o /xray

FROM alpine:3.9
# hadolint ignore=DL3022
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /xray /usr/bin/
RUN ln -sf /dev/stdout /var/log/xray-daemon.log
EXPOSE 2000/udp
COPY xray-daemon.yaml /
ENTRYPOINT ["/usr/bin/xray", "--config", "/xray-daemon.yaml"]
CMD ["--region", "us-west-2"]
