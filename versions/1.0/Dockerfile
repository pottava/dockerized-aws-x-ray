FROM alpine:3

RUN apk --no-cache add ca-certificates
RUN apk --no-cache add --virtual build-dependencies bash curl libgcc unzip \
    && GLIBC_VERSION=2.29-r0 \
    && GLIBC_DL_URL="https://github.com/andyshinn/alpine-pkg-glibc/releases/download/${GLIBC_VERSION}" \
    && curl --location --silent --show-error -O ${GLIBC_DL_URL}/glibc-${GLIBC_VERSION}.apk \
    && curl --location --silent --show-error -O ${GLIBC_DL_URL}/glibc-bin-${GLIBC_VERSION}.apk \
    && curl --location --silent --show-error -O ${GLIBC_DL_URL}/glibc-i18n-${GLIBC_VERSION}.apk \
    && apk add --allow-untrusted glibc-${GLIBC_VERSION}.apk \
       glibc-bin-${GLIBC_VERSION}.apk glibc-i18n-${GLIBC_VERSION}.apk \
    && /usr/glibc-compat/sbin/ldconfig /lib /usr/glibc-compat/lib \
    && /usr/glibc-compat/bin/localedef -i en_US -f UTF-8 en_US.UTF-8 \
    && rm -rf glibc-${GLIBC_VERSION}.apk \
       glibc-bin-${GLIBC_VERSION}.apk glibc-i18n-${GLIBC_VERSION}.apk \
    && cd /usr/lib \
    && curl -O --location --silent --show-error https://s3.amazonaws.com/aws-xray-assets.us-east-1/xray-daemon/aws-xray-daemon-linux-1.x.zip \
    && unzip -q aws-xray-daemon-linux-1.x.zip \
    && rm -rf aws-xray-daemon-linux-1.x.zip \
    && ln -s /usr/lib/xray /usr/bin/ \
    && ln -sf /dev/stdout /var/log/xray-daemon.log \
    && apk del --purge -r build-dependencies

EXPOSE 2000

ENTRYPOINT ["xray"]
CMD ["--bind", "0.0.0.0:2000", "--log-file", "/var/log/xray-daemon.log", "--memory-limit", "100"]
