# AWS X-Ray v3.2
#
# docker run --rm pottava/xray:3.2-arm --help
# docker run --name xray -d -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY \
#     -p 2000:2000/udp pottava/xray:3.2-arm --region ${AWS_REGION} \
#     --local-mode --log-level debug

FROM alpine:3.10 AS certs
RUN apk --no-cache add ca-certificates

FROM golang:1.13.5-alpine3.10 AS build
RUN apk --no-cache add git
RUN go get -u github.com/aws/aws-xray-daemon/...
WORKDIR /go/src/github.com/aws/aws-xray-daemon
ENV XRAY_VERSION=3.2.0
RUN git checkout "${XRAY_VERSION}"
WORKDIR /go/src/github.com/aws/aws-xray-daemon/daemon
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm64
RUN go build -ldflags "-s -w" -o /xray
RUN ln -sf /dev/stdout /var/log/xray-daemon.log

FROM scratch
EXPOSE 2000/tcp
EXPOSE 2000/udp
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /var/log /var/log
COPY --from=build /xray /usr/bin/
COPY xray-daemon.yaml /
ENTRYPOINT ["/usr/bin/xray", "--config", "/xray-daemon.yaml"]
CMD ["--region", "us-west-2"]
